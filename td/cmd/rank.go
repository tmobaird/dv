package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

var RankCmd = &cobra.Command{
	Use:   "rank [from] [to]",
	Short: "Moves an existing todo from one position to another.",
	Long:  "Moves an existing todo from one position to another. For example, if you want to make a task more important, you can move it higher up. If it's less important, you can move it lower down.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.RankController{Base: controllers.Controller{Args: args, Config: config}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
