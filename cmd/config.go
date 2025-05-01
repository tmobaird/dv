package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"td/internal"

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
		dirname := ".td"
		if os.Getenv("TD_BASE_PATH") != "" {
			dirname = os.Getenv("TD_BASE_PATH")
		}
		if Edit {
			if os.Getenv("EDITOR") == "" {
				cmd.OutOrStderr().Write([]byte("Must set $EDITOR to edit config"))
				return
			}
			cmd := exec.Command(os.Getenv("EDITOR"), fmt.Sprintf("%s/config.yaml", dirname)) // Replace filename.txt with the actual file path
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			config, err := internal.Read(os.DirFS(dirname))
			if err != nil {
				cmd.OutOrStderr().Write([]byte("FAILED TO READ"))
				cmd.OutOrStderr().Write([]byte(err.Error()))
			} else {
				cmd.OutOrStderr().Write([]byte(fmt.Sprintf("Current Config:\n  Context: %s", config.Context)))
			}
		}
	},
}
