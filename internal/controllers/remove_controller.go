package controllers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	todos, err := controller.getTodos()
	if err != nil {
		return "", err
	}
	if todoIndex < 1 || (todoIndex-1) > len(todos) {
		return "", fmt.Errorf("todo @ index %d does not exist", todoIndex)
	}
	target := todos[todoIndex-1]
	mdTodos := controller.writableTodos(todos, todoIndex)
	err = controller.writeTodos(mdTodos)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\"%s\" removed from list.", target.Name), nil
}

func (controller RemoveController) getTodos() ([]models.Todo, error) {
	todos, err := ReadTodos(controller.Base.Config.Context)
	if err != nil {
		return []models.Todo{}, err
	}
	filtered := FilterTodos(todos, controller.Base.Config.HideCompleted)
	return filtered, err
}

func (controller RemoveController) writableTodos(todos []models.Todo, index int) []string {
	toWrite := []string{}
	for i, todo := range todos {
		if i+1 != index {
			toWrite = append(toWrite, models.TodoToMd(todo))
		}
	}
	return toWrite
}

func (controller RemoveController) writeTodos(mdTodos []string) error {
	file, err := os.OpenFile(internal.TodoFilePath(controller.Base.Config.Context), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content := strings.Join(mdTodos, "\n")
	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}
