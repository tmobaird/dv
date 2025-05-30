package main

import (
	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/td"
)

var rootCmd = &cobra.Command{
	Use:   "dv",
	Short: "dv is a developer tools cli",
	Long:  "replace me",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(td.TdCommand())
	rootCmd.Execute()
}
