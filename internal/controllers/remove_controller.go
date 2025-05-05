package controllers

import (
	"fmt"
	"strconv"
	"td/internal"
	"td/internal/models"
)

type RemoveController struct {
	Base Controller
}

func (controller RemoveController) Run() (string, error) {
	todoArg, err := strconv.Atoi(controller.Base.Args[0])
	if err != nil {
		return "", err
	}
	todoIndex := todoArg - 1

	todos, err := models.GetAllTodos(internal.TodoFilePath(controller.Base.Config.Context))
	if err != nil {
		return "", err
	}
	if todoIndex < 0 || todoIndex >= len(todos) {
		return "", fmt.Errorf("todo @ index %d does not exist", todoArg)
	}

	if controller.Base.Config.HideCompleted {
		notCompletedCounter := 0
		for i, todo := range todos {
			if !todo.Complete {
				if notCompletedCounter == todoIndex {
					todoIndex = i
					break
				}
				notCompletedCounter++
			}
		}
	}

	target := todos[todoIndex]
	todos = controller.writableTodos(todos, todoIndex)
	err = models.WriteTodos(controller.Base.Config.Context, todos)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\"%s\" removed from list.", target.Name), nil
}

func (controller RemoveController) writableTodos(todos []models.Todo, index int) []models.Todo {
	toWrite := []models.Todo{}
	for i, todo := range todos {
		if i != index {
			toWrite = append(toWrite, todo)
		}
	}
	return toWrite
}
