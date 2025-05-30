package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

var RenameCmd = &cobra.Command{
	Use:     "rename [index] [new-name]",
	Short:   "Renames an existing item in the todo list",
	Long:    "Renames an existing todo in the list in your current context. Simply pass the index of the todo and a string to this command to update the name.",
	Aliases: []string{"mv"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.RenameController{Base: controllers.Controller{Args: args, Config: config}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
