package controllers

import (
	"context"
	"errors"
	"io"
	"os"
	"td/internal"
	"td/internal/auth"
	"td/internal/models"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type ScheduleController struct {
	Base       Controller
	Regenerate bool
	NoCalendar bool
}

func (controller ScheduleController) Run() (string, error) {
	d := time.Now()

	if controller.Regenerate {
		result, err := generate(d, !controller.NoCalendar)
		return result, err
	} else {
		result, err := readSchedule(d)
		if err != nil {
			result, err = generate(d, !controller.NoCalendar)
			return result, err
		}
		return result, err
	}
}

func generate(d time.Time, useCalendar bool) (string, error) {
	cal := models.Calendar{}

	if useCalendar {
		ctx := context.Background()
		b, err := os.ReadFile("credentials.json") // Downloaded from Google Cloud Console
		if err != nil {
			return "", err
		}
		config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			return "", err
		}
		client := auth.GetClient(ctx, config)
		cal, err = models.GetTodaysCalendar(client, ctx)
		if err != nil {
			return "", errors.New("failed to fetch calendar")
		}
	}

	schedule, err := models.CreateSchedule(d, cal, "main")
	return schedule, err
}

func readSchedule(d time.Time) (string, error) {
	file, err := os.OpenFile(internal.ScheduleFilePath(d), os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
