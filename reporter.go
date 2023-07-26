package main

import "fmt"

func ReportTodos(todos []Todo, verbose bool) string {
	var output string

	for index, todo := range todos {
		output += fmt.Sprintf("%d. ", index+1)
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
