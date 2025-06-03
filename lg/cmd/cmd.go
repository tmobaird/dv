package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/lg/internal/controllers"
)

var WriteCmd = &cobra.Command{
	Use:     "write",
	Aliases: []string{"w"},
	Short:   "writes some shit",
	Long:    "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.Controller{Args: args, Config: config}.RunWrite()

		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "reads the shit",
	Long:  "TODO",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		result, err := controllers.Controller{Args: args, Config: config}.RunShow()

		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}

func Commands() []*cobra.Command {
	return []*cobra.Command{WriteCmd, ShowCmd}
}
