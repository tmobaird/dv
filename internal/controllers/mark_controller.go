package controllers

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"td/internal"
	"td/internal/models"
)

type MarkController struct {
	Base Controller
}

const DONE = "done"
const D = "d"
const NOT_DONE = "not-done"
const NOT = "not"

func (controller MarkController) Run() (string, error) {
	todoNum, err := strconv.Atoi(controller.Base.Args[0])
	if err != nil {
		return "", err
	}
	todoIndex := todoNum - 1

	todos, err := models.GetAllTodos(internal.TodoFilePath(controller.Base.Config.Context))
	if err != nil {
		return "", err
	}

	err = controller.validateArgs(todoIndex, controller.Base.Args[1], todos)
	if err != nil {
		return "", err
	}

	// update single todo
	todoIndex = controller.findCorrectIndex(todoIndex, todos)
	todos[todoIndex].Complete = controller.markCompleted(controller.Base.Args[1])
	target := todos[todoIndex]

	// write todos
	err = models.WriteTodos(controller.Base.Config.Context, todos)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\" marked %s.", target.Name, target.Status()), nil
}

func (controller MarkController) validateArgs(index int, status string, todos []models.Todo) error {
	if index < 0 || index >= len(todos) {
		return fmt.Errorf("todo @ index %d does not exist", index+1)
	}
	allowed := []string{DONE, D, NOT_DONE, NOT}
	if !slices.Contains(allowed, status) {
		return fmt.Errorf("status must be one of: [%s], got %s", strings.Join(allowed, ", "), status)
	}
	return nil
}

func (controller MarkController) findCorrectIndex(initial int, todos []models.Todo) int {
	if controller.Base.Config.HideCompleted {
		notCompletedCounter := 0
		for i, todo := range todos {
			if !todo.Complete {
				if notCompletedCounter == initial {
					return i
				}
				notCompletedCounter++
			}
		}
	}
	return initial
}

func (controller MarkController) markCompleted(status string) bool {
	return status == DONE || status == D
}
