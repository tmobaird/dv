package cmd

import (
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	ScheduleCmd.Flags().BoolVarP(&Regenerate, "regenerate", "r", false, "Regenerate")
	rootCmd.AddCommand(ScheduleCmd)
}

var Regenerate bool
var ScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Show the current schedule for today.",
	Long:  "By default shows the current config. When --edit is passed, it will open the config file in your current $EDITOR.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := controllers.ScheduleController{Base: controllers.Controller{Args: args}, Regenerate: Regenerate}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
