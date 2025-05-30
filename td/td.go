package td

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/td/cmd"
)

var tdCmd = &cobra.Command{
	Use:   "td",
	Short: "td is a todo list manager",
	Long:  "TD is a simple todo list manager that supports multiple lists, todo statuses, and AI-based scheduling.\nThis tool allows users to manage and optimize their time.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func TdCommand() *cobra.Command {
	commands := cmd.Commands()
	for _, command := range commands {
		tdCmd.AddCommand(command)
	}
	return tdCmd
}

func Execute() {
	// if err := core.LoadConfig(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	if err := tdCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
