package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"td/internal/auth"
	"td/internal/models"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type CalendarController struct {
	Base Controller
}

func (controller CalendarController) Run() (string, error) {
	ctx := context.Background()

	// Load OAuth 2.0 config from client_secrets.json
	b, err := os.ReadFile("credentials.json") // Downloaded from Google Cloud Console
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Set the desired scopes
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := auth.GetClient(ctx, config)

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
