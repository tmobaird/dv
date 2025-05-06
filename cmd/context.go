package cmd

import (
	"os"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	ContextCmd.Flags().BoolVarP(&Default, "default", "d", false, "Go back to default context (ie main)")
	rootCmd.AddCommand(ContextCmd)
}

var Default bool
var ContextCmd = &cobra.Command{
	Use:     "context [name?]",
	Short:   "Get or set the current context",
	Long:    "Get or set the current context. When name is given, the current context will be updated.",
	Aliases: []string{"pwd"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.ContextController{Base: controllers.Controller{Args: args, Config: config}, ChangeToDefault: Default}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
