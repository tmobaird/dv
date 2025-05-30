package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
)

var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open current todo list for editing",
	Long:  "Open current todo list for editing. Relies on $EDITOR to open the current todo list in your editor of choice. The list will be written in markdown.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("EDITOR") == "" {
			cmd.OutOrStderr().Write([]byte("Must set $EDITOR to edit config"))
			return
		}
		config, err := core.ReadConfig(os.DirFS(core.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		openEditorCmd := exec.Command(os.Getenv("EDITOR"), core.TodoFilePath(config.Context))
		openEditorCmd.Stdin = os.Stdin
		openEditorCmd.Stdout = os.Stdout
		err = openEditorCmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}
