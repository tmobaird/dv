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
	t.Run("#generateSchedule can generate a simple schedule", func(t *testing.T) {
		todoContent := "- [x] Todo A\n- [ ] Todo B\n- [ ] Todo C"
		dirname := CreateTodoSetup(t, "main", []byte(todoContent))
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfDay(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 2, len(schedule.Items))
	})

	t.Run("#generateSchedule will prioritize deep work in the first 1/3 of day", func(t *testing.T) {
		deep := Todo{Name: "Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}
		shallow := Todo{Name: "Shallow", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Shallow}}
		admin := Todo{Name: "Admin", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Blank}}
		todos := []Todo{shallow, admin, deep}

		contents := []string{}
		for _, todo := range todos {
			contents = append(contents, todo.ToMd())
		}

		dirname := CreateTodoSetup(t, "main", []byte(strings.Join(contents, "")))
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfDay(time.Now()), []CalendarEvent{}, "main")

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
		todos := []Todo{deep, deep, deep, deep, deep, deep, deep, shallow, shallow, shallow, shallow, keyDeep}

		contents := []string{}
		for _, todo := range todos {
			contents = append(contents, todo.ToMd())
		}

		dirname := CreateTodoSetup(t, "main", []byte(strings.Join(contents, "")))
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(startOfWorkday(time.Now()), []CalendarEvent{}, "main")

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 12, len(schedule.Items))
		testutils.AssertEqual(t, "Key Deep", schedule.Items[len(schedule.Items)-1].Item.Output())
	})

	t.Run("#generateSchedule will not schedule todos that conflict with meetings", func(t *testing.T) {
		todo := Todo{Name: "Deep", Complete: false, Metadata: Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Deep}}

		dirname := CreateTodoSetup(t, "main", []byte(todo.ToMd()))
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

		dirname := CreateTodoSetup(t, "main", []byte(todo.ToMd()))
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

	t.Run("#generateSchedule can handle back to back meetings", func(t *testing.T) {
		dirname := CreateTodoSetup(t, "main", []byte(""))
		defer os.RemoveAll(dirname)

		schedule, err := generateSchedule(time.Now(), []CalendarEvent{
			{StartTime: time.Now(), EndTime: time.Now().Add(time.Minute * 30), Event: &calendar.Event{Summary: "Meeting A"}},
			{StartTime: time.Now().Add(time.Minute * 30), EndTime: time.Now().Add(time.Minute * 60), Event: &calendar.Event{Summary: "Meeting B"}},
		}, "main")

		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, 2, len(schedule.Items))
	})

	t.Run("#generateSchedule can correctly handle a schedule with no tasks, and only meetings", func(t *testing.T) {
		dirname := CreateTodoSetup(t, "main", []byte(""))
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
