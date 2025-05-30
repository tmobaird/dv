package controllers

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/auth"
	"github.com/tmobaird/dv/td/internal/models"
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
		client, err := auth.GetClient(ctx)
		if err != nil {
			return "", err
		}

		cal, err = models.GetTodaysCalendar(client, ctx)
		if err != nil {
			return "", errors.New("failed to fetch calendar")
		}
	}

	schedule, err := models.CreateSchedule(d, cal, "main")
	return schedule, err
}

func readSchedule(d time.Time) (string, error) {
	file, err := os.OpenFile(core.ScheduleFilePath(d), os.O_RDWR, 0644)
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
