package main

import (
	"bytes"
	"os"
	"strings"
	"td/cmd"
	"td/internal"
	"testing"

	"github.com/spf13/cobra"
)

func TestIntegration(t *testing.T) {
	t.Run("Config returns config", func(t *testing.T) {
		// setup
		os.Setenv("TD_BASE_PATH", "tmp")
		file, err := os.Create("tmp/config.yaml")
		if err != nil {
			t.Errorf("Failed to create tmp file %s", err.Error())
		}
		defer os.Remove(file.Name())
		config := internal.Config{Context: "main"}
		internal.Save(file, config)

		// execution
		rootCmd := &cobra.Command{Use: "app"}
		rootCmd.AddCommand(cmd.ConfigCmd)
		outputBuf := new(bytes.Buffer)
		rootCmd.SetOut(outputBuf)
		rootCmd.SetArgs([]string{"config"})
		if err := rootCmd.Execute(); err != nil {
			t.Errorf("command failed: %v", err)
		}

		// expectation
		got := outputBuf.String()
		if !strings.Contains(got, "Context: main") {
			t.Errorf("Expected %s to contain %s", got, "Context: main")
		}
	})

	t.Run("Context", func(t *testing.T) {

	})
}
