package main

import (
	"fmt"
	"strings"
)

func ReportTodos(todos []Todo, verbose bool) string {
	var output string

	for index, todo := range todos {
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
