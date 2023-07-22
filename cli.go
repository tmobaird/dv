package main

import "fmt"

type Cli struct {
	Command string
	Args []string
	Commander Commander
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
	}

	if err != nil {
		fmt.Printf("ERROR - Failed to run command: %s", err)
	} else {
		fmt.Println("Todos")
		fmt.Println("=====")
		for _, todo := range todos {
			fmt.Printf("%s - %s\n", todo.Id, todo.Name)
		}
	}
}