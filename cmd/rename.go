package cmd

import (
	"os"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RenameCmd)
}

var RenameCmd = &cobra.Command{
	Use:     "rename [index] [new-name]",
	Aliases: []string{"mv"},
	Short:   "Renames an existing item in the todo list",
	Long:    `Renames an existing todo in the list in your current context. Simply pass the index of the todo and a string to this command to update the name.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
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
