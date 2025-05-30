package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/tmobaird/dv/core"
)

const DEEP_WORK_PRIORITIZED_IN = float64(0.33)

type ScheduledItem interface {
	Output() string
}

type ScheduleItem struct {
	Start time.Time
	End   time.Time
	Item  ScheduledItem
}

type Schedule struct {
	At    time.Time
	Items []ScheduleItem
}

func (schedule Schedule) Print() string {
	result := fmt.Sprintf("Here is your schedule for %s\n", schedule.At.Format(time.DateOnly))
	for _, item := range schedule.Items {
		result += fmt.Sprintf("%s-%s: %s\n", item.Start.Format(time.Kitchen), item.End.Format(time.Kitchen), item.Item.Output())
	}
	return result
}

func CreateSchedule(date time.Time, calendar Calendar, context string) (string, error) {
	err := os.MkdirAll(core.ScheduleDirectoryPath(), 0700)
	if err != nil {
		return "", err
	}
	schedule, err := generateSchedule(date, calendar.Events, context)
	if err != nil {
		return "", err
	}
	contents := schedule.Print()
	file, err := os.OpenFile(core.ScheduleFilePath(date), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(contents))
	if err != nil {
		return "", err
	}
	// return contents
	return contents, nil
}

func generateSchedule(start time.Time, events []CalendarEvent, context string) (Schedule, error) {
	initialStart := start
	dayDuration := endOfWorkday(initialStart).Sub(initialStart)
	todos, err := GetAllTodos(core.TodoFilePath(context))
	if err != nil {
		return Schedule{}, err
	}
	todos = FilterTodos(todos, true)
	todoBank := map[int]time.Duration{}
	for i, todo := range todos {
		todoBank[i] = todo.Duration()
	}

	schedule := Schedule{Items: []ScheduleItem{}, At: time.Now()}
	for (len(availableTodos(todos, todoBank)) > 0 || len(events) > 0) && start.Before(endOfWorkday(time.Now())) {
		availableTime := time.Duration(-1)

		var upcomingEvent CalendarEvent
		if len(events) > 0 {
			upcomingEvent = events[0]
			availableTime = upcomingEvent.StartTime.Sub(start)
		}

		duration := start.Sub(initialStart)
		percentage := (float64(duration) / float64(dayDuration))
		prioritizeDeepWork := percentage < float64(DEEP_WORK_PRIORITIZED_IN)
		todo, index, err := findTodo(todos, todoBank, availableTime, prioritizeDeepWork)

		if err != nil {
			schedule.Items = append(schedule.Items, scheduleItemForEvent(upcomingEvent))
			events = events[1:]
			start = upcomingEvent.EndTime
		} else {
			updateBankDuration := todoBank[index]
			addUpcomingEvent := false
			if len(events) > 0 && start.Add(todo.Duration()).After(upcomingEvent.StartTime) {
				updateBankDuration = upcomingEvent.StartTime.Sub(start)
				addUpcomingEvent = true
			}

			schedule.Items = append(schedule.Items, scheduleItemForTodo(todo, start, updateBankDuration))
			if addUpcomingEvent {
				schedule.Items = append(schedule.Items, scheduleItemForEvent(upcomingEvent))
				events = events[1:]
				start = upcomingEvent.EndTime
			} else {
				start = start.Add(updateBankDuration)
			}
			updateBank(index, updateBankDuration, todoBank)
		}
	}

	for len(events) > 0 {
		schedule.Items = append(schedule.Items, scheduleItemForEvent(events[0]))
		events = events[1:]
	}

	return schedule, nil
}

func updateBank(index int, duration time.Duration, bank map[int]time.Duration) {
	bank[index] -= duration
}

func availableTodos(todos []Todo, bank map[int]time.Duration) []Todo {
	available := []Todo{}
	for i, todo := range todos {
		if bank[i] > 0 {
			available = append(available, todo)
		}
	}

	return available
}

func findTodo(todos []Todo, todoBank map[int]time.Duration, maxDuration time.Duration, prioritizeDeepWork bool) (Todo, int, error) {
	if prioritizeDeepWork {
		for i, todo := range todos {
			if todo.Metadata.Type == Deep &&
				todoBank[i] > 0 &&
				(maxDuration == -1 || maxDuration >= todo.Duration()) {
				return todo, i, nil
			}
		}
	}

	for i, todo := range todos {
		if todoBank[i] > 0 {
			return todo, i, nil
		}
	}

	return Todo{}, -1, errors.New("failed to find a todo")
}

func scheduleItemForTodo(todo Todo, start time.Time, duration time.Duration) ScheduleItem {
	return ScheduleItem{Item: todo, Start: start, End: start.Add(duration)}
}

func scheduleItemForEvent(event CalendarEvent) ScheduleItem {
	return ScheduleItem{Item: event, Start: event.StartTime, End: event.EndTime}
}
