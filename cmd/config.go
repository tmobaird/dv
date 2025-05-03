package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"td/internal"
	"td/internal/controllers"

	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.Flags().BoolVarP(&Edit, "edit", "e", false, "Edit config")
	rootCmd.AddCommand(ConfigCmd)
}

var Edit bool
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get current config",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if Edit {
			if os.Getenv("EDITOR") == "" {
				cmd.OutOrStderr().Write([]byte("Must set $EDITOR to edit config"))
				return
			}
			cmd := exec.Command(os.Getenv("EDITOR"), internal.ConfigFilePath()) // Replace filename.txt with the actual file path
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
