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

	t.Run("adds and returns multiple todos at once", func(t *testing.T) {
		todosGot, _ := add([]string{"Hello", "World"}, &MockReaderWriter{})
		todosWant := []Todo{
			{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()},
			{Name: "World", CreatedAt: "dummy", Id: uuid.New()},
		}
		got := len(todosGot)
		want := len(todosWant)

		if got != want {
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

		got, _ := delete([]string{todoUuid.String()}, rw)
		want := 1

		assertLength(t, got, want)
	})

	t.Run("deletes multiple todos from list", func(t *testing.T) {
		todoUuid := uuid.New()
		todoUuidTwo := uuid.New()
		todoOne := Todo{Name: "Hello", CreatedAt: "dummy", Id: todoUuid}
		todoTwo := Todo{Name: "World", CreatedAt: "dummy", Id: todoUuidTwo}
		rw := &MockReaderWriter{todos: []Todo{todoOne, todoTwo}}

		got, _ := delete([]string{todoUuid.String(), todoUuidTwo.String()}, rw)
		want := 0

		assertLength(t, got, want)
	})

	t.Run("can delete by index", func(t *testing.T) {
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()}}}

		got, _ := delete([]string{"1"}, rw)
		want := 0

		assertLength(t, got, want)
	})

	t.Run("returns error when fails to read file", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := delete([]string{uuid.New().String()}, rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})

	t.Run("returns error when fails to write file", func(t *testing.T) {
		rw := &ErrorMockWriter{}

		_, gotError := delete([]string{uuid.New().String()}, rw)
		wantErrorMessage := "Failed to write file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})
}

func TestDone(t *testing.T) {
	t.Run("marks todo as done", func(t *testing.T) {
		todoUuid := uuid.New()
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: todoUuid}, {Name: "World", CreatedAt: "dummy", Id: uuid.New()}}}

		todos, _ := done(todoUuid.String(), rw)
		got := todos[0]
		want := true


		assertDone(t, got, want)
	})

	t.Run("marks todo as done by index", func(t *testing.T) {
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: uuid.New(), Done: false}}}

		todos, _ := done("1", rw)
		got := todos[0]
		want := true

		assertDone(t, got, want)
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
		got := todos[0]
		want := false

		assertDone(t, got, want)
	})

	t.Run("marks todo as not done by index", func(t *testing.T) {
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: uuid.New(), Done: true}}}

		todos, _ := undo("1", rw)
		got := todos[0]
		want := false

		assertDone(t, got, want)
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

func TestEdit(t *testing.T) {
	assertString := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got \"%s\" want \"%s\"", got, want)
		}
	}

	t.Run("edits name of todo", func(t *testing.T) {
		todoUuid := uuid.New()
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: todoUuid}}}
		newName := "new name"

		todos, _ := edit(todoUuid.String(), newName, rw)
		got := todos[0].Name
		want := newName


		assertString(t, got, want)
	})

	t.Run("edits name of todo by index", func(t *testing.T) {
		rw := &MockReaderWriter{todos: []Todo{{Name: "Hello", CreatedAt: "dummy", Id: uuid.New()}}}
		newName := "new name"

		todos, _ := edit("1", newName, rw)
		got := todos[0].Name
		want := newName


		assertString(t, got, want)
	})

	t.Run("returns error when fails to read file", func(t *testing.T) {
		rw := &ErrorMockReader{}

		_, gotError := edit(uuid.New().String(), "new name", rw)
		wantErrorMessage := "Failed to read file"

		assertCorrectError(t, gotError, wantErrorMessage)
	})

	t.Run("returns error when fails to write file", func(t *testing.T) {
		rw := &ErrorMockWriter{}

		_, gotError := edit(uuid.New().String(), "new name", rw)
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

func assertLength(t testing.TB, got []Todo, want int) {
	t.Helper()
	if len(got) != want {
		t.Errorf("got %v of len %d want %v", got, len(got), want)
	}
}

func assertDone(t testing.TB, got Todo, want bool) {
	t.Helper()
	if got.Done != want {
		t.Errorf("got %v want %v", got.Done, want)
	}
}