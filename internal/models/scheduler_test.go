package models

import (
	"fmt"
	"os"
	"strings"
	"td/internal/testutils"
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
)

func TestScheduler(t *testing.T) {
	t.Run("#generateSchedule schedules todos based on priority", func(t *testing.T) {
		a := Todo{Name: "Todo A", Complete: true, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		b := Todo{Name: "Todo B", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		c := Todo{Name: "Todo C", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}

		dirname := CreateTodosSetup(t, "main", []Todo{a, b, c})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfDay(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 2, len(schedule.Items))
		testutils.AssertEqual(t, "Todo B", schedule.Items[0].Item.Output())
		testutils.AssertEqual(t, "Todo C", schedule.Items[1].Item.Output())
	})

	t.Run("#generateSchedule Schedules deep work todos first in the beginning of day", func(t *testing.T) {
		deep := Todo{Name: "Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		shallow := Todo{Name: "Shallow", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Shallow}}
		admin := Todo{Name: "Admin", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Blank}}

		dirname := CreateTodosSetup(t, "main", []Todo{shallow, admin, deep})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 3, len(schedule.Items))
		testutils.AssertEqual(t, "Deep", schedule.Items[0].Item.Output())
		testutils.AssertEqual(t, "Shallow", schedule.Items[1].Item.Output())
		testutils.AssertEqual(t, "Admin", schedule.Items[2].Item.Output())
	})

	t.Run("#generateSchedule will NOT prioritize deep work in the latter 2/3 of day", func(t *testing.T) {
		deep := Todo{Name: "Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		keyDeep := Todo{Name: "Key Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		shallow := Todo{Name: "Shallow", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Shallow}}

		dirname := CreateTodosSetup(t, "main", []Todo{deep, deep, deep, deep, deep, deep, deep, shallow, shallow, shallow, shallow, keyDeep})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 12, len(schedule.Items))
		testutils.AssertEqual(t, "Key Deep", schedule.Items[len(schedule.Items)-1].Item.Output())
	})

	t.Run("#generateSchedule schedules tasks for the appropriate amount of time", func(t *testing.T) {
		a := Todo{Name: "Todo A", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		b := Todo{Name: "Todo B", Complete: false, Metadata: Metadata{DurationValue: 15, DurationUnit: MINUTE_UNIT, Type: Deep}}
		c := Todo{Name: "Todo C", Complete: false, Metadata: Metadata{DurationValue: 2, DurationUnit: HOUR_UNIT, Type: Deep}}

		dirname := CreateTodosSetup(t, "main", []Todo{a, b, c})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{}, "main")
		format := "2006-01-02 15:04"
		var offset time.Duration = 0

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 3, len(schedule.Items))

		testutils.AssertEqual(t, startOfWorkday(time.Now()).Add(offset).Format(format), schedule.Items[0].Start.Format(format))
		testutils.AssertEqual(t, startOfWorkday(time.Now()).Add(offset).Add(unitToDuration(BLOCK_UNIT, 1)).Format(format), schedule.Items[0].End.Format(format))
		offset += unitToDuration(BLOCK_UNIT, 1)

		testutils.AssertEqual(t, startOfWorkday(time.Now()).Add(offset).Format(format), schedule.Items[1].Start.Format(format))
		testutils.AssertEqual(t, startOfWorkday(time.Now()).Add(offset).Add(unitToDuration(MINUTE_UNIT, 15)).Format(format), schedule.Items[1].End.Format(format))
		offset += unitToDuration(MINUTE_UNIT, 15)

		testutils.AssertEqual(t, startOfWorkday(time.Now()).Add(offset).Format(format), schedule.Items[2].Start.Format(format))
		testutils.AssertEqual(t, startOfWorkday(time.Now()).Add(offset).Add(unitToDuration(HOUR_UNIT, 2)).Format(format), schedule.Items[2].End.Format(format))
	})

	t.Run("#generateSchedule will not schedule todos that conflict with meetings", func(t *testing.T) {
		todo := Todo{Name: "Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}

		dirname := CreateTodosSetup(t, "main", []Todo{todo})
		defer os.RemoveAll(dirname)

		meetingStartTime := startOfWorkday(time.Now()).Add(5 * time.Minute)
		meetingEndTime := startOfWorkday(time.Now()).Add(35 * time.Minute)
		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{
			{Event: &calendar.Event{Summary: "Meeting A"}, StartTime: meetingStartTime, EndTime: meetingEndTime},
		}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 3, len(schedule.Items))
		testutils.AssertEqual(t, "Deep", schedule.Items[0].Item.Output())
		testutils.AssertEqual(t, meetingStartTime, schedule.Items[0].End)
		testutils.AssertEqual(t, "Meeting A", schedule.Items[1].Item.Output())
		testutils.AssertEqual(t, "Deep", schedule.Items[len(schedule.Items)-1].Item.Output())
		testutils.AssertEqual(t, meetingEndTime, schedule.Items[len(schedule.Items)-1].Start)
	})

	t.Run("#generateSchedule will gracefully handle durations for todos that are broken up by meetings", func(t *testing.T) {
		todo := Todo{Name: "Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}

		dirname := CreateTodosSetup(t, "main", []Todo{todo})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{
			{Event: &calendar.Event{Summary: "Meeting A"}, StartTime: startOfWorkday(time.Now()).Add(5 * time.Minute), EndTime: startOfWorkday(time.Now()).Add(35 * time.Minute)},
		}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 3, len(schedule.Items))
		testutils.AssertEqual(t, "Deep", schedule.Items[0].Item.Output())
		testutils.AssertEqual(t, "Meeting A", schedule.Items[1].Item.Output())
		testutils.AssertEqual(t, "Deep", schedule.Items[len(schedule.Items)-1].Item.Output())
	})

	t.Run("#generateSchedule does not schedule completed todos", func(t *testing.T) {
		todo := Todo{Name: "Todo A", Complete: true, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}

		dirname := CreateTodosSetup(t, "main", []Todo{todo})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfDay(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 0, len(schedule.Items))
	})

	t.Run("#generateSchedule when no tasks exist and only calendar events, it adds the events to the calendar", func(t *testing.T) {
		dirname := CreateTodosSetup(t, "main", []Todo{})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(time.Now(), []CalendarEvent{
			{StartTime: time.Now(), EndTime: time.Now().Add(time.Minute * 30), Event: &calendar.Event{Summary: "Meeting A"}},
			{StartTime: time.Now().Add(time.Minute * 30), EndTime: time.Now().Add(time.Minute * 60), Event: &calendar.Event{Summary: "Meeting B"}},
		}, "main")

		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, 2, len(schedule.Items))
		testutils.AssertEqual(t, "Meeting A", schedule.Items[0].Item.Output())
		testutils.AssertEqual(t, "Meeting B", schedule.Items[1].Item.Output())
	})

	t.Run("#generateSchedule back to back meetings are scheduled gracefully", func(t *testing.T) {
		dirname := CreateTodosSetup(t, "main", []Todo{})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(time.Now(), []CalendarEvent{
			{StartTime: time.Now(), EndTime: time.Now().Add(time.Minute * 30), Event: &calendar.Event{Summary: "Meeting A"}},
			{StartTime: time.Now().Add(time.Minute * 30), EndTime: time.Now().Add(time.Minute * 60), Event: &calendar.Event{Summary: "Meeting B"}},
		}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 2, len(schedule.Items))
	})

	t.Run("#generateSchedule only schedules til 6:30 PM", func(t *testing.T) {
		a := Todo{Name: "Todo A", Complete: false, Metadata: Metadata{DurationValue: 10, DurationUnit: HOUR_UNIT, Type: Deep}}
		b := Todo{Name: "Todo B", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}

		dirname := CreateTodosSetup(t, "main", []Todo{a, b})
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 1, len(schedule.Items))
		testutils.AssertEqual(t, "Todo A", schedule.Items[0].Item.Output())
	})
}

func CreateTodosSetup(t *testing.T, context string, todos []Todo) string {
	t.Helper()

	contents := []string{}
	for _, todo := range todos {
		contents = append(contents, todo.ToMd())
	}

	dirname := CreateTodoSetup(t, "main", []byte(strings.Join(contents, "")))
	return dirname
}

func CreateTodoSetup(t *testing.T, context string, contents []byte) string {
	t.Helper()

	os.Setenv("TD_BASE_PATH", "tmp")
	dirname := "tmp/lists"
	err := os.MkdirAll(dirname, 0755)
	testutils.AssertNoError(t, err)
	file, err := os.Create(fmt.Sprintf("tmp/lists/%s.md", context))
	testutils.AssertNoError(t, err)
	file.Write(contents)
	return dirname
}
