package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

var MarkCmd = &cobra.Command{
	Use:     "mark [index] [done|d|not-done|not]",
	Short:   "Update a todo's current status",
	Long:    "Update a todo's current status (done, d, not-done, not)",
	Aliases: []string{"m"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.MarkController{Base: controllers.Controller{Args: args, Config: config}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
