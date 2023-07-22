package main

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestReportTodos(t *testing.T) {
	t.Run("turns todos into string output", func(t *testing.T) {
		todos := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()},
			{Name: "World", CreatedAt: "dummy", Id: uuid.New()},
		}

		got := ReportTodos(todos)
		want := "- Hello\n- World\n"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestReportError(t *testing.T) {
	t.Run("turns error into string output", func(t *testing.T) {
		err := errors.New("Failed to write to file")

		got := ReportError(err, "add")
		want := "ERROR - Failed to execute add\n- Failed to write to file\n"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}