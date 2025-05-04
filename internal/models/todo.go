package models

import "fmt"

type Todo struct {
	Name     string
	Complete bool
}

func TodoToMd(todo Todo) string {
	completeChar := " "
	if todo.Complete {
		completeChar = "x"
	}
	return fmt.Sprintf("- [%s] %s\n", completeChar, todo.Name)
}

