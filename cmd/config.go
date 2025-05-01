package cmd

import (
	"fmt"
	"os"
	"td/internal"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.Flags().BoolVarP(&Edit, "edit", "e", false, "Edit config")
	rootCmd.AddCommand(configCmd)
}

var Edit bool
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get current config",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if Edit {
			fmt.Println(Edit)
		} else {
			config, err := internal.Read(os.DirFS(".td"))
			if err != nil {
				fmt.Println("FAILED TO READ CONFIG", err.Error())
			} else {
				fmt.Println("GOT CONFIG", config)
			}
		}
	},
}
