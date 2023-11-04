package main

import (
	"fmt"
	"reflect"
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

func ReportConfig(config Config) string {
	output := ""

	output += "=== Config ===\n"
	v := reflect.ValueOf(config)
    t := v.Type()

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        output += fmt.Sprintf("%s: %v\n", field.Name, value.Interface())
    }

	return output
}
