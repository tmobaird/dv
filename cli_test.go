package main

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("calls add command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockAdd := func(_ [] string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "add",
			Args:      []string{"Hello", "World"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls add command with a", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockAdd := func(_ [] string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "a",
			Args:      []string{"Hello", "World"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls command with args", func(t *testing.T) {
		var gotArgs []string
		wantArgs := []string{"Hello", "World"}
		mockAdd := func(args [] string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			gotArgs = args
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "add",
			Args:      []string{"Hello", "World"},
			Commander: commander,
		}

		cli.Run()

		if reflect.DeepEqual(gotArgs, wantArgs) == false {
			t.Errorf("Expected %v, got %v", wantArgs, gotArgs)
		}
	})

	t.Run("does not call command when name doesnt match", func(t *testing.T) {
		count := 0
		wantCount := 0
		mockAdd := func(args [] string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Add: mockAdd, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "not-echo",
			Args:      []string{"Hello", "World"},
			Commander: commander,
		}

		cli.Run()

		if count != wantCount {
			t.Errorf("Expected %d, got %d", wantCount, count)
		}
	})

	t.Run("calls list command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockList := func(r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{List: mockList, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "list",
			Commander: commander,
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
		commander := Commander{List: mockList, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "ls",
			Commander: commander,
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
		commander := Commander{Delete: mockDelete, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "delete",
			Args: 	[]string{"123"},
			Commander: commander,
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
		commander := Commander{Delete: mockDelete, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "d",
			Args: 	[]string{"123"},
			Commander: commander,
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
		commander := Commander{Done: mockDone, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "done",
			Args: 	[]string{"123"},
			Commander: commander,
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
		commander := Commander{Done: mockDone, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "do",
			Args: 	[]string{"123"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls undo command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockUndo := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Undo: mockUndo, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "undo",
			Args: 	[]string{"123"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls undo command with un", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockUndo := func(_ string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Undo: mockUndo, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "un",
			Args: 	[]string{"123"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls edit command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockEdit := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Edit: mockEdit, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "edit",
			Args: 	[]string{"123", "new name"},
			Commander: commander,
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
		commander := Commander{Edit: mockEdit, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "e",
			Args: 	[]string{"123", "new name"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})

	t.Run("calls rank command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockRank := func(_, _b string, r ReaderWriter) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commander := Commander{Rank: mockRank, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "rank",
			Args: 	[]string{"123", "1"},
			Commander: commander,
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
		commander := Commander{Rank: mockRank, ReaderWriter: &RealReaderWriter{}}

		cli := Cli{
			Command:   "r",
			Args: 	[]string{"123", "1"},
			Commander: commander,
		}

		cli.Run()

		assertCalls(t, count, wantCount)
	})
}


func assertCalls(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("Expected %d, got %d", want, got)
	}
}