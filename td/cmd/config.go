package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/controllers"
)

var Edit bool
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows the current config.",
	Long:  "By default shows the current config. When --edit is passed, it will open the config file in your current $EDITOR.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if Edit {
			if os.Getenv("EDITOR") == "" {
				cmd.OutOrStderr().Write([]byte("Must set $EDITOR to edit config"))
				return
			}
			cmd := exec.Command(os.Getenv("EDITOR"), core.ConfigFilePath()) // Replace filename.txt with the actual file path
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			result, err := controllers.ConfigController{Base: controllers.Controller{Args: args}}.Run()
			if err != nil {
				cmd.OutOrStderr().Write([]byte(err.Error()))
			} else {
				cmd.OutOrStdout().Write([]byte(result))
			}
		}
	},
}
