package models

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"td/internal"
)

type Todo struct {
	Name     string
	Complete bool
}

func (todo *Todo) ToMd() string {
	completeChar := " "
	if todo.Complete {
		completeChar = "x"
	}
	return fmt.Sprintf("- [%s] %s\n", completeChar, todo.Name)
}

func GetAllTodos(filename string, hideCompleted bool) ([]Todo, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return []Todo{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return []Todo{}, err
	}

	todos := parseTodos(data)
	todos = filterTodos(todos, hideCompleted)
	return todos, nil
}

func WriteTodos(context string, todos []Todo) error {
	file, err := os.OpenFile(internal.TodoFilePath(context), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content := ""
	for _, todo := range todos {
		content += fmt.Sprintf("%s\n", todo.ToMd())
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}

func parseTodos(bytes []byte) []Todo {
	todos := []Todo{}
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

				todos = append(todos, Todo{Name: parts[2], Complete: complete})
			}
		}
	}
	return todos
}

func filterTodos(todos []Todo, hideCompleted bool) []Todo {
	result := []Todo{}
	for _, todo := range todos {
		if hideCompleted && todo.Complete {
			// do nothing, ie filter out
		} else {
			result = append(result, todo)
		}
	}
	return result
}
