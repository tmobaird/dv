package cmd

import (
	"os"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RemoveCmd)
}

var RemoveCmd = &cobra.Command{
	Use:     "remove [index]",
	Short:   "Remove a todo from list",
	Long:    "Removes a todo from the list in your current context. Simply pass the index of the todo in the list.",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
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
