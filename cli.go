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

	switch {
		case c.Command == "list":
			todos, err = c.Commander.List(c.Commander.ReaderWriter)
		case c.Command == "delete":
			todos, err = c.Commander.Delete(c.Args[0], c.Commander.ReaderWriter)
		case c.Command == "add":
			todos, err = c.Commander.Add(c.Args, c.Commander.ReaderWriter)
		case c.Command == "done":
			todos, err = c.Commander.Done(c.Args[0], c.Commander.ReaderWriter)
		case c.Command == "undo":
			todos, err = c.Commander.Undo(c.Args[0], c.Commander.ReaderWriter)
		default:
			fmt.Println("No command provided")
	}

	if err != nil {
		fmt.Printf(ReportError(err, c.Command))
	} else {
		fmt.Print(ReportTodos(todos, c.Verbose))
	}
}