package cmd

import (
	"fmt"
	"os"
	"td/internal"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ContextCmd)
}

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Get or set the current context",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.Read(os.DirFS(internal.BasePath()))
		if err != nil {
			cmd.OutOrStderr().Write([]byte(err.Error()))
			return
		}

		if len(args) > 0 {
			config.Context = args[0]
			err := internal.PersistConfig(config)
			if err != nil {
				cmd.OutOrStderr().Write([]byte(err.Error()))
			} else {
				cmd.OutOrStdout().Write([]byte(fmt.Sprintf("Updated context to %s", config.Context)))
			}
		} else {
			cmd.OutOrStdout().Write([]byte(config.Context))
		}
	},
}
