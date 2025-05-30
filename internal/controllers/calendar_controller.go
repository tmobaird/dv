package controllers

import (
	"context"
	"fmt"
	"td/internal/auth"
	"td/internal/models"
	"time"
)

type CalendarController struct {
	Base Controller
}

func (controller CalendarController) Run() (string, error) {
	ctx := context.Background()
	client, err := auth.GetClient(ctx)
	if err != nil {
		return "", err
	}

	calendar, err := models.GetTodaysCalendar(client, ctx)
	if err != nil {
		return "", err
	}

	result := "Here's a look at your calendar for today:\n"
	for i, item := range calendar.Events {
		result += fmt.Sprintf("%d. %s - %s: %s\n", i+1, item.StartTime.Format(time.TimeOnly), item.EndTime.Format(time.TimeOnly), item.Event.Summary)
	}

	return result, nil
}
