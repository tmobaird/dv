package controllers

import (
	"fmt"
	"strconv"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/models"
)

type RenameController struct {
	Base Controller
}

func (controller RenameController) Run() (string, error) {
	// get target index
	todoNum, err := strconv.Atoi(controller.Base.Args[0])
	if err != nil {
		return "", err
	}
	todoIndex := todoNum - 1

	// get todos
	todos, err := models.GetAllTodos(core.TodoFilePath(controller.Base.Config.Context))
	if err != nil {
		return "", err
	}

	// error when out of range
	if todoIndex < 0 || todoIndex >= len(todos) {
		return "", fmt.Errorf("todo @ index %d does not exist", todoNum)
	}

	if controller.Base.Config.HideCompleted {
		notCompletedCounter := 0
		startSet := false
		for i, todo := range todos {
			if !todo.Complete {
				if notCompletedCounter == todoIndex && !startSet {
					todoIndex = i
					startSet = true
				}
				notCompletedCounter++
			}
		}
	}

	// find todo and update todo
	todo := todos[todoIndex]
	previousName := todo.Name
	todos[todoIndex].Name = controller.Base.Args[1]

	// save todos
	err = models.WriteTodos(controller.Base.Config.Context, todos)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("\"%s\" updated to \"%s\".", previousName, todos[todoIndex].Name), nil
	}
}
