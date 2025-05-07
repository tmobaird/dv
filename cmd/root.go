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
	Long:  "TD is a simple todo list manager that supports multiple lists, todo statuses, and AI-based scheduling.\nThis tool allows users to manage and optimize their time.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func loadConfig() error {
	dirname := internal.BasePath()
	if os.Getenv("TD_BASE_PATH") != "" {
		dirname = os.Getenv("TD_BASE_PATH")
	}
	config, err := internal.Read(os.DirFS(dirname))
	if err != nil {
		return err
	}

	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}

	err = internal.PersistConfig(config)
	if err != nil {
		return err
	}
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
