package main

import (
	"fmt"
	"strings"
)

func completedTodos(todos []Todo) int {
	var completed int
	for _, todo := range todos {
		if todo.Done {
			completed++
		}
	}
	return completed
}

func ReportTodos(todos []Todo, verbose, hideCompleted bool) string {
	var output string

	numCompleted := completedTodos(todos)
	if hideCompleted && numCompleted > 0 {
		output += fmt.Sprintf("=== Hidden: %d completed todos ===\n", numCompleted)
	}

	for index, todo := range todos {
		if hideCompleted && todo.Done {
			continue
		}
		trailingSpaces := 2
		if index+1 >= 10 {
			trailingSpaces = 1
		}
		output += fmt.Sprintf("%d.%s", index+1, (strings.Repeat(" ", trailingSpaces)))
		if todo.Done {
			output += "[✓] "
		} else {
			output += "[ ] "
		}
		output += todo.Name

		if verbose {
			output += " (" + todo.Id.String() + ")"
		}
		output += "\n"
	}

	return output
}

func ReportError(error error, command string) string {
	return fmt.Sprintf("ERROR - Failed to execute %s\n- %s\n", command, error.Error())
}
