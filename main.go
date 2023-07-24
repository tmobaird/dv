package main

import (
	"fmt"
	"os"
)

func helpIfNecessary(args []string) {
	help := false

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			help = true
		}
	}

	if help {
		fmt.Println("Usage: td [options] [command] [arguments]")
		fmt.Println("Options:")
		fmt.Println("  -h, --help     Print usage")
		fmt.Println("  -v, --verbose  Print verbose output")
		fmt.Println("Commands:")
		fmt.Println("  add <name>              Add a new todo")
		fmt.Println("  list                    List all todos")
		fmt.Println("  delete <uuid>           Delete a todo")
		fmt.Println("  done <uuid>             Mark a todo as done")
		fmt.Println("  undo <uuid>             gMark a todo as not done")
		fmt.Println("  edit <uuid> <new name>  Edit a todo")
		os.Exit(0)
	}
}

func main() {
	args := os.Args[1:]
	verbose := false

	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		}
	}

	helpIfNecessary(os.Args[1:])

	if len(args) > 0 {
		cmd := args[0]
		commander := Commander{Add: add, List: list, Delete: delete, Done: done, Undo: undo, Edit: edit, ReaderWriter: &RealReaderWriter{}}
		cli := Cli{Command: cmd, Args: args[1:], Commander: commander, Verbose: verbose}
		cli.Run()
	} else {
		fmt.Println("No command provided")
		fmt.Println("Print Help")
	}
}
