package controllers

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
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

	filename := fmt.Sprintf("%s/lists/%s.md", internal.BasePath(), controller.Base.Config.Context)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	todos := []models.Todo{}
	parseTodos(data, &todos)

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

func parseTodos(bytes []byte, todos *[]models.Todo) {
	data := string(bytes)
	initialRegex := regexp.MustCompile("^- .*")              // checks to make sure it is a valid list item
	partsRegex := regexp.MustCompile(`^- \[([^\]]+)\] (.*)`) // captures complete and name
	for _, line := range strings.Split(data, "\n") {
		if initialRegex.Match([]byte(line)) {
			parts := partsRegex.FindStringSubmatch(line)
			complete := false
			if len(parts) >= 3 {
				if parts[1] == "x" {
					complete = true
				}

				*todos = append(*todos, models.Todo{Name: parts[2], Complete: complete})
			}
		}
	}
}

func presentTodo(index int, todo models.Todo) string {
	completechar := " "
	if todo.Complete {
		completechar = "x"
	}

	return fmt.Sprintf("%d. [%s] %s\n", index, completechar, todo.Name)
}
