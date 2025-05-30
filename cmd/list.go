package cmd

import (
	"os"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	ListCmd.Flags().BoolVarP(&All, "all", "a", false, "Show all")
	ListCmd.Flags().BoolVar(&ShowMetadata, "metadata", false, "Show todo metadata")
	rootCmd.AddCommand(ListCmd)
}

var All bool
var ShowMetadata bool
var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List todos for current context",
	Long:    "List todos for current context. Pass -a or --all to show completed todos.",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		if All {
			config.HideCompleted = false
		}
		result, err := controllers.ListController{Base: controllers.Controller{Args: args, Config: config}, ShowMetadata: ShowMetadata}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
