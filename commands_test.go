package main

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

type MockReaderWriter struct {
	todos []Todo
}

func (m *MockReaderWriter) ReadJSONFileToMap() ([]Todo, error) {
	return m.todos, nil
}

func (m *MockReaderWriter) WriteTodosToFile(todos []Todo) error {
	m.todos = todos
	return nil
}

type ErrorMockReader struct{}

func (m *ErrorMockReader) ReadJSONFileToMap() ([]Todo, error) {
	return nil, errors.New("Failed to read file")
}

func (m *ErrorMockReader) WriteTodosToFile(todos []Todo) error {
	return nil
}

type ErrorMockWriter struct{}

func (m *ErrorMockWriter) ReadJSONFileToMap() ([]Todo, error) {
	return []Todo{}, nil
}

func (m *ErrorMockWriter) WriteTodosToFile(todos []Todo) error {
	return errors.New("Failed to write file")
}

func TestAdd(t *testing.T) {
	t.Run("returns todos with new items", func(t *testing.T) {
		got, _ := add([]string{"Hello"}, &MockReaderWriter{})
		want := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()},
		}

		if got[0].Name != want[0].Name {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("writes todos to file", func(t *testing.T) {
		rw := &MockReaderWriter{}
		add([]string{"Hello"}, rw)

		got, _ := rw.ReadJSONFileToMap()
		want := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()},
		}

		if got[0].Name != want[0].Name {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns error when file cannot be read", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := add([]string{"Hello"}, rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})

	t.Run("returns error when file cannot be written", func(t *testing.T) {
		rw := &ErrorMockWriter{}

		_, gotError := add([]string{"Hello"}, rw)
		wantErrorMessage := "Failed to write file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})
}

func TestList(t *testing.T) {
	t.Run("lists todos", func(t *testing.T) {
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()}}}
		got, _ := list(rw)
		want := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()},
		}

		if got[0].Name != want[0].Name {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns error when file cannot be read", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := list(rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})
}

func TestDelete(t *testing.T) {
	t.Run("deletes todo from list", func(t *testing.T) {
		todoUuid := uuid.New()
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: todoUuid}, {Name: "World", CreatedAt: "dummy", Id: uuid.New()}}}

		got, _ := delete(todoUuid.String(), rw)
		want := 1

		if len(got) != want {
			t.Errorf("got %v of len %d want %v", got, len(got), want)
		}
	})

	t.Run("returns error when fails to read file", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := delete(uuid.New().String(), rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})

	t.Run("returns error when fails to write file", func(t *testing.T) {
		rw := &ErrorMockWriter{}

		_, gotError := delete(uuid.New().String(), rw)
		wantErrorMessage := "Failed to write file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})
}

func TestDone(t *testing.T) {
	t.Run("marks todo as done", func(t *testing.T) {
		todoUuid := uuid.New()
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: todoUuid}, {Name: "World", CreatedAt: "dummy", Id: uuid.New()}}}

		todos, _ := done(todoUuid.String(), rw)
		got := todos[0].Done
		want := true


		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns error when fails to read file", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := done(uuid.New().String(), rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})

	t.Run("returns error when fails to write file", func(t *testing.T) {
		rw := &ErrorMockWriter{}

		_, gotError := done(uuid.New().String(), rw)
		wantErrorMessage := "Failed to write file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})
}

func TestUndo(t *testing.T) {
	t.Run("marks todo as not done", func(t *testing.T) {
		todoUuid := uuid.New()
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: todoUuid, Done: true}}}

		todos, _ := undo(todoUuid.String(), rw)
		got := todos[0].Done
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns error when fails to read file", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := undo(uuid.New().String(), rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})

	t.Run("returns error when fails to write file", func(t *testing.T) {
		rw := &ErrorMockWriter{}

		_, gotError := undo(uuid.New().String(), rw)
		wantErrorMessage := "Failed to write file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})
}

func assertCorrectError(t testing.TB, got error, want string) {
	t.Helper()
	if got.Error() != want {
		t.Errorf("got \"%s\" want \"%s\"", got, want)
	}
}