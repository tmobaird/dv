package main

import "fmt"

func ReportTodos (todos []Todo) string {
	var output string

	for _, todo := range todos {
		output += "- " + todo.Name + "\n"
	}

	return output
}

func ReportError (error error, command string) string {
	return fmt.Sprintf("ERROR - Failed to execute %s\n- %s\n", command, error.Error())
}