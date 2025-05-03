package cmd

import (
	"os"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ContextCmd)
}

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Get or set the current context",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.ContextController{Base: controllers.Controller{Args: args, Config: config}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
