package models

import (
	"fmt"
	"os"
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
