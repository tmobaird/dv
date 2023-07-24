package main

import (
	"fmt"
)

type Cli struct {
	Command string
	Args []string
	Commander Commander
	Verbose bool
}

func (c Cli) Run() {
	var todos []Todo
	var err error

	switch c.Command {
		case "list":
			todos, err = c.Commander.List(c.Commander.ReaderWriter)
		case "delete":
			todos, err = c.Commander.Delete(c.Args[0:], c.Commander.ReaderWriter)
		case "add":
			todos, err = c.Commander.Add(c.Args, c.Commander.ReaderWriter)
		case "done":
			todos, err = c.Commander.Done(c.Args[0], c.Commander.ReaderWriter)
		case "undo":
			todos, err = c.Commander.Undo(c.Args[0], c.Commander.ReaderWriter)
		case "edit":
			todos, err = c.Commander.Edit(c.Args[0], c.Args[1], c.Commander.ReaderWriter)
		default:
			fmt.Println("No command provided")
	}

	if err != nil {
		fmt.Printf(ReportError(err, c.Command))
	} else {
		fmt.Print(ReportTodos(todos, c.Verbose))
	}
}