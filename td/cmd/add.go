package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

func init() {
	AddCmd.Flags().StringVarP(&Duration, "duration", "d", "1b", "Estimated duration of todo")
	AddCmd.Flags().StringVarP(&Type, "type", "t", "blank", "Set the type of todo")
	// rootCmd.AddCommand(AddCmd)
}

var Duration string
var Type string
var AddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new todo to list",
	Long:  `Adds a new todo to the list in your current context. Simply pass a string to this command to add the todo.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		metadataString := fmt.Sprintf("duration=%s,type=%s", Duration, Type)
		result, err := controllers.AddController{Base: controllers.Controller{Args: args, Config: config}, MetadataString: metadataString}.Run()
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
		} else {
			cmd.OutOrStdout().Write([]byte(result))
		}
	},
}

func AddCommand() *cobra.Command {
	return AddCmd
}
