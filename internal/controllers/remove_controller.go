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
	todoIndex, err := strconv.Atoi(controller.Base.Args[0])
	if err != nil {
		return "", err
	}

	todos, err := models.GetAllTodos(internal.TodoFilePath(controller.Base.Config.Context), controller.Base.Config.HideCompleted)
	if err != nil {
		return "", err
	}
	if todoIndex < 1 || (todoIndex-1) > len(todos) {
		return "", fmt.Errorf("todo @ index %d does not exist", todoIndex)
	}

	target := todos[todoIndex-1]
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
		if i+1 != index {
			toWrite = append(toWrite, todo)
		}
	}
	return toWrite
}
