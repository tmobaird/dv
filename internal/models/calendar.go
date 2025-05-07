package models

import (
	"context"
	"net/http"
	"sort"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarEvent struct {
	Event     *calendar.Event
	StartTime time.Time
	EndTime   time.Time
}

type Calendar struct {
	Events []CalendarEvent
	Start  time.Time
	End    time.Time
}

func GetTodaysCalendar(client *http.Client, ctx context.Context) (Calendar, error) {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return Calendar{}, err
	}

	timeMin := startOfDay(time.Now())
	timeMax := endOfDay(time.Now())
	events, err := srv.Events.List("primary").TimeMin(timeMin.Format(time.RFC3339)).TimeMax(timeMax.Format(time.RFC3339)).Do()
	if err != nil {
		return Calendar{}, err
	}

	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].Start.DateTime < events.Items[j].Start.DateTime // Sort by age
	})

	calendarEvents := []CalendarEvent{}
	for _, event := range events.Items {
		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)
		calendarEvents = append(calendarEvents, CalendarEvent{Event: event, StartTime: start, EndTime: end})
	}

	return Calendar{Events: calendarEvents, Start: timeMin, End: timeMax}, nil
}

func startOfDay(t time.Time) time.Time {
	location, _ := time.LoadLocation("Local")
	year, month, day := t.In(location).Date()
	t = time.Date(year, month, day, 0, 0, 0, 0, location)
	return t
}

func endOfDay(t time.Time) time.Time {
	location, _ := time.LoadLocation("Local")
	year, month, day := t.In(location).Date()
	t = time.Date(year, month, day, 23, 59, 59, 0, location)
	return t
}
