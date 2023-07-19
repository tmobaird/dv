package main

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("calls command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockAdd := func([] string) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commandController := CommandController{add: mockAdd}

		cli := Cli{
			Command:   "add",
			Args:      []string{"Hello", "World"},
			Commander: commandController,
		}

		cli.Run()

		if count != wantCount {
			t.Errorf("Expected %d, got %d", wantCount, count)
		}
	})

	t.Run("calls command with args", func(t *testing.T) {
		var gotArgs []string
		wantArgs := []string{"Hello", "World"}
		mockAdd := func(args [] string) ([]Todo, error) {
			todos := []Todo{}
			gotArgs = args
			return todos, nil
		}
		commandController := CommandController{add: mockAdd}

		cli := Cli{
			Command:   "add",
			Args:      []string{"Hello", "World"},
			Commander: commandController,
		}

		cli.Run()

		if reflect.DeepEqual(gotArgs, wantArgs) == false {
			t.Errorf("Expected %v, got %v", wantArgs, gotArgs)
		}
	})

	t.Run("does not call command when name doesnt match", func(t *testing.T) {
		count := 0
		wantCount := 0
		mockAdd := func([] string) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commandController := CommandController{add: mockAdd}

		cli := Cli{
			Command:   "not-echo",
			Args:      []string{"Hello", "World"},
			Commander: commandController,
		}

		cli.Run()

		if count != wantCount {
			t.Errorf("Expected %d, got %d", wantCount, count)
		}
	})

	t.Run("calls list command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockList := func() ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commandController := CommandController{list: mockList}

		cli := Cli{
			Command:   "list",
			Commander: commandController,
		}

		cli.Run()

		if count != wantCount {
			t.Errorf("Expected %d, got %d", wantCount, count)
		}
	})

	t.Run("calls delete command", func(t *testing.T) {
		count := 0
		wantCount := 1
		mockDelete := func(_ string) ([]Todo, error) {
			todos := []Todo{}
			count += 1
			return todos, nil
		}
		commandController := CommandController{delete: mockDelete}

		cli := Cli{
			Command:   "delete",
			Args: 	[]string{"123"},
			Commander: commandController,
		}

		cli.Run()

		if count != wantCount {
			t.Errorf("Expected %d, got %d", wantCount, count)
		}
	})
}
