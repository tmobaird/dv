package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("calls add command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockAdd := func(_ []string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "add",
			Args:      []string{"Hello", "World"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls add command with a", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockAdd := func(_ []string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "a",
			Args:      []string{"Hello", "World"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls command with args", func(t *testing.T) {
		var gotArgs []string
		wantArgs := []string{"Hello", "World"}
		mockAdd := func(args []string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			gotArgs = args
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "add",
			Args:      []string{"Hello", "World"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		if reflect.DeepEqual(gotArgs, wantArgs) == false {
			t.Errorf("Expected %v, got %v", wantArgs, gotArgs)
		}
	})

	t.Run("does not call command when name doesnt match", func(t *testing.T) {
		count := 0
		mockAdd := func(args []string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "not-echo",
			Args:      []string{"Hello", "World"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, 0)
	})

	t.Run("calls list command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockList := func(r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{List: mockList, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "list",
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls list command with ls", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockList := func(r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{List: mockList, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "ls",
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls delete command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockDelete := func(_ []string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Delete: mockDelete, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "delete",
			Args:      []string{"123"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls delete command with d", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockDelete := func(_ []string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Delete: mockDelete, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "d",
			Args:      []string{"123"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls done command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockDone := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Done: mockDone, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "done",
			Args:      []string{"123"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls done command with do", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockDone := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Done: mockDone, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "do",
			Args:      []string{"123"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("done raises error when less than 1 arg", func(t *testing.T) {
		count := 0
		data := []string{}
		mockDone := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Done: mockDone, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "done",
			Args:      []string{},
			Commander: commander,
			PrintFunc: func(a ...interface{}) (n int, err error) {
				data = append(data, fmt.Sprint(a...))
				return 0, nil
			},
		}

		cli.Run()

		assertCalls(t, count, 0)
		assertErrorIncludingMessage(t, data[0], NoTodoSpecifedErr)
	})

	t.Run("calls undo command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockUndo := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Undo: mockUndo, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "undo",
			Args:      []string{"123"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls undo command with un", func(t *testing.T) {
		count := 0
		wantCount := 1
		data := []string{}
		mockUndo := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Undo: mockUndo, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "un",
			Args:      []string{"123"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
		fmt.Println(data)
	})

	t.Run("undo raises error when less than 1 arg", func(t *testing.T) {
		count := 0
		data := []string{}
		mockUndo := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Undo: mockUndo, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "un",
			Args:      []string{},
			Commander: commander,
			PrintFunc: func(a ...interface{}) (n int, err error) {
				data = append(data, fmt.Sprint(a...))
				return 0, nil
			},
		}

		cli.Run()

		assertCalls(t, count, 0)
		assertErrorIncludingMessage(t, data[0], NoTodoSpecifedErr)
	})

	t.Run("calls edit command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockEdit := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Edit: mockEdit, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "edit",
			Args:      []string{"123", "new name"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls edit with e", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockEdit := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Edit: mockEdit, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "e",
			Args:      []string{"123", "new name"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("edit raises error when less than 2 args", func(t *testing.T) {
		count := 0
		data := []string{}
		mockEdit := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Edit: mockEdit, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "e",
			Args:      []string{},
			Commander: commander,
			PrintFunc: func(a ...interface{}) (n int, err error) {
				data = append(data, fmt.Sprint(a...))
				return 0, nil
			},
		}

		cli.Run()

		assertCalls(t, count, 0)
		assertErrorIncludingMessage(t, data[0], NoIndexOrNameErr)
	})

	t.Run("calls rank command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockRank := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Rank: mockRank, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "rank",
			Args:      []string{"123", "1"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls rank command with r", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockRank := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Rank: mockRank, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "r",
			Args:      []string{"123", "1"},
			Commander: commander,
			PrintFunc: fmt.Print,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("rank raises error when less than 2 args", func(t *testing.T) {
		count := 0
		data := []string{}
		mockRank := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Rank: mockRank, ReaderWriter: &TdReaderWriter{}}

		cli := Cli{
			Command:   "rank",
			Args:      []string{},
			Commander: commander,
			PrintFunc: func(a ...interface{}) (n int, err error) {
				data = append(data, fmt.Sprint(a...))
				return 0, nil
			},
		}

		cli.Run()

		assertCalls(t, count, 0)
		assertErrorIncludingMessage(t, data[0], NoIndexOrRankErr)
	})
}

func assertCalls(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("Expected %d, got %d", want, got)
	}
}

func assertErrorIncludingMessage(t testing.TB, got string, want error) {
	t.Helper()

	if strings.Contains(got, want.Error()) == false {
		t.Errorf("Expected \"%s\", got \"%s\"", want.Error(), got)
	}
}
