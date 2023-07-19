package main

import "fmt"

type Cli struct {
	Command string
	Args []string
	Commander CommandController
}

func (c Cli) Run() {
	var todos []Todo
	var err error

	switch {
		case c.Command == "list":
			todos, err = c.Commander.list()
		case c.Command == "delete":
			todos, err = c.Commander.delete(c.Args[0])
		case c.Command == "add":
			todos, err = c.Commander.add(c.Args)
	}

	if err != nil {
		fmt.Printf("Failed to run command: %s", err)
	} else {
		for _, todo := range todos {
			fmt.Printf("%s - %s\n", todo.Id, todo.Name)
		}
	}
}