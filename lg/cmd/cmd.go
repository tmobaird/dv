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
	Short:   "opens a file for writing your daily developer log",
	Long:    "write: opens a new file in .dv/logs/<today>.md from the set template to allow you to write your daily developer log. Set $EDITOR to change which program the file opens in for editing.",
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
			cmd.OutOrStdout().Write([]byte(result.String()))
		}
	},
}

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "shows the latest developer log you have completed",
	Long:  "Show: shows the latest developer log you have completed",
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
			cmd.OutOrStdout().Write([]byte(result.String()))
		}
	},
}

var Before string
var After string
var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "shows you a list of your past developer logs",
	Long:  "log: similar to git log, this shows you your list of past developer logs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		logArgs := controllers.LogArgs{Before: Before, After: After}
		result, err := controllers.Controller{Args: args, Config: config}.RunLog(logArgs)

		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result.String()))
		}
	},
}

func init() {
	LogCmd.Flags().StringVar(&Before, "before", "", "Before date")
	LogCmd.Flags().StringVar(&Before, "until", "", "Before date")
	LogCmd.Flags().StringVar(&After, "after", "", "After date")
	LogCmd.Flags().StringVar(&After, "since", "", "After date")
}

func Commands() []*cobra.Command {
	return []*cobra.Command{WriteCmd, ShowCmd, LogCmd}
}
