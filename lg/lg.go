package lg

import (
	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/lg/cmd"
)

var lgCmd = &cobra.Command{
	Use:   "lg",
	Short: "lg is a developer log system",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func LgCommand() *cobra.Command {
	commands := cmd.Commands()
	for _, command := range commands {
		lgCmd.AddCommand(command)
	}
	return lgCmd
}
