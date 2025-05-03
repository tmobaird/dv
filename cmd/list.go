package cmd

import (
	"os"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	ListCmd.Flags().BoolVarP(&Edit, "edit", "e", false, "Edit config")
	rootCmd.AddCommand(ListCmd)
}

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List todos for space",
	Long:    `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.ListController{Base: controllers.Controller{Args: args, Config: config}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
