package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		cmd := args[0]
		// commander := CommandController{add: add, list: list, delete: delete}
		commander := Commander{Add: add, List: list, Delete: delete, ReaderWriter: &RealReaderWriter{}}
		cli := Cli{Command: cmd, Args: args[1:], Commander: commander}
		cli.Run()
	} else {
		fmt.Println("No command provided")
		fmt.Println("Print Help")
	}
}