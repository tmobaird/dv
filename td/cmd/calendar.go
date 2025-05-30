package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

var CalendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Add a new todo to list",
	Long:  `Adds a new todo to the list in your current context. Simply pass a string to this command to add the todo.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}
		result, err := controllers.CalendarController{Base: controllers.Controller{Config: config, Args: args}}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}
