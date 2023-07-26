package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("Usage: td [options] [command] [arguments]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help     Print usage")
	fmt.Println("  -v, --verbose  Print verbose output")
	fmt.Println("Commands:")
	fmt.Println("  a,  add <name>                      Add a new todo")
	fmt.Println("  ls, list                            List all todos")
	fmt.Println("  d,  delete <index|uuid>             Delete a todo")
	fmt.Println("  do, done   <index|uuid>             Mark a todo as done")
	fmt.Println("  un, undo   <index|uuid>             Mark a todo as not done")
	fmt.Println("  e,  edit   <index|uuid> <new name>  Edit a todo")
	fmt.Println("  r,  rank   <index|uuid> <new rank>  Rerank a todo")
	os.Exit(0)
}

func helpIfNecessary(args []string) {
	needHelp := false

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			needHelp = true
		}
	}

	if needHelp {
		help()
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
		commander := Commander{
			Add:          add,
			List:         list,
			Delete:       delete,
			Done:         done,
			Undo:         undo,
			Edit:         edit,
			Rank:         rank,
			ReaderWriter: &RealReaderWriter{},
		}
		cli := Cli{Command: cmd, Args: args[1:], Commander: commander, Verbose: verbose}
		cli.Run()
	} else {
		fmt.Println("No command provided")
		help()
	}
}
