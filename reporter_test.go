package main

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestReportTodos(t *testing.T) {
	t.Run("turns todos into string output", func(t *testing.T) {
		todos := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: uuid.New(), Done: false},
			{Name: "World", CreatedAt: "dummy", Id: uuid.New(), Done: true},
		}

		got := ReportTodos(todos, false, false)
		want := "1.  [ ] Hello\n2.  [✓] World\n"

		assertCorrectMessage(t, got, want)
	})

	t.Run("shows todo id when verbose", func(t *testing.T) {
		todoId := uuid.New()
		todos := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: todoId, Done: false},
		}

		got := ReportTodos(todos, true, false)
		want := "1.  [ ] Hello (" + todos[0].Id.String() + ")\n"

		assertCorrectMessage(t, got, want)
	})

	t.Run("hides completed todos when hideCompleted is true", func(t *testing.T) {
		todoId := uuid.New()
		todos := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: todoId, Done: true},
		}

		got := ReportTodos(todos, true, true)
		want := "=== Hidden: 1 completed todos ===\n"

		assertCorrectMessage(t, got, want)
	})

	t.Run("prints nothing when there are no todos", func(t *testing.T) {
		todos := []Todo{}

		got := ReportTodos(todos, true, true)
		want := ""

		assertCorrectMessage(t, got, want)
	})
}

func TestReportError(t *testing.T) {
	t.Run("turns error into string output", func(t *testing.T) {
		err := errors.New("Failed to write to file")

		got := ReportError(err, "add")
		want := "ERROR - Failed to execute add\n- Failed to write to file\n"

		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
