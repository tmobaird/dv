package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

var RemoveCmd = &cobra.Command{
	Use:     "remove [index]",
	Short:   "Remove a todo from list",
	Long:    "Removes a todo from the list in your current context. Simply pass the index of the todo in the list.",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.RemoveController{Base: controllers.Controller{Args: args, Config: config}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
