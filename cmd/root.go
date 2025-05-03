package cmd

import (
	"fmt"
	"os"
	"td/internal"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "td",
	Short: "td is a todo list manager",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func loadConfig() error {
	dirname := ".td"
	if os.Getenv("TD_BASE_PATH") != "" {
		dirname = os.Getenv("TD_BASE_PATH")
	}
	config, err := internal.Read(os.DirFS(dirname))

	if err != nil {
		return err
	}

	file, err := os.OpenFile(dirname, os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	internal.Save(file, config)
	return nil
}

func Execute() {
	if err := loadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
