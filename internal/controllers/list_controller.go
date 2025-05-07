package controllers

import (
	"fmt"
	"os"
	"td/internal"
	"td/internal/models"
)

type ListController struct {
	Base Controller
}

func (controller ListController) Run() (string, error) {
	dirPath := fmt.Sprintf("%s/lists", internal.BasePath())
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", err
	}

	todos, err := models.GetAllTodos(internal.TodoFilePath(controller.Base.Config.Context))
	todos = models.FilterTodos(todos, controller.Base.Config.HideCompleted)
	if err != nil {
		return "", err
	}

	if len(todos) > 0 {
		result := ""
		for i, todo := range todos {
			result += presentTodo(i+1, todo)
		}

		return result, nil
	} else {
		return "No todos in list.", nil
	}
}

func presentTodo(index int, todo models.Todo) string {
	completechar := " "
	if todo.Complete {
		completechar = "x"
	}

	return fmt.Sprintf("%d. [%s] %s\n", index, completechar, todo.Name)
}
