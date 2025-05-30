package controllers

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/models"
)

type RankController struct {
	Base Controller
}

func (controller RankController) Run() (string, error) {
	parsedArgs, err := controller.parseArgs()
	if err != nil {
		return "", err
	}
	todoArg := parsedArgs[0]
	endArg := parsedArgs[1]
	todoIndex := todoArg - 1
	endIndex := endArg - 1

	todos, err := models.GetAllTodos(core.TodoFilePath(controller.Base.Config.Context))
	if err != nil {
		return "", err
	}
	err = controller.validateArgs(todoIndex, endIndex, todos)
	if err != nil {
		return "", err
	}

	if controller.Base.Config.HideCompleted {
		notCompletedCounter := 0
		startSet := false
		endSet := false
		for i, todo := range todos {
			if !todo.Complete {
				if notCompletedCounter == todoIndex && !startSet {
					todoIndex = i
					startSet = true
				}
				if notCompletedCounter == endIndex && !endSet {
					endIndex = i
					endSet = true
				}
				notCompletedCounter++
			}
		}
	}

	// find, delete from list, insert into list
	target := todos[todoIndex]
	todos = slices.Delete(todos, todoIndex, todoIndex+1)
	todos = slices.Insert(todos, endIndex, target)

	err = models.WriteTodos(controller.Base.Config.Context, todos)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\"%s\" ranked to index %d.", target.Name, endArg), nil
}

func (controller RankController) parseArgs() ([]int, error) {
	todoArg, todoArgErr := strconv.Atoi(controller.Base.Args[0])
	endArg, endArgErr := strconv.Atoi(controller.Base.Args[1])
	if todoArgErr != nil || endArgErr != nil {
		return []int{}, fmt.Errorf("rank inputs are invalid: %s, %s", controller.Base.Args[0], controller.Base.Args[1])
	}
	return []int{todoArg, endArg}, nil
}

func (controller RankController) validateArgs(startIndex, endIndex int, todos []models.Todo) error {
	if startIndex < 0 || startIndex > len(todos) {
		return fmt.Errorf("todo @ index %d does not exist", startIndex+1)
	}
	if endIndex < 0 || endIndex > len(todos) {
		return fmt.Errorf("cant move todo to @%d, out of bounds", endIndex+1)
	}
	return nil
}
