package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func CreateConfigFile(t *testing.T) *os.File {
	t.Helper()

	os.Setenv("TD_BASE_PATH", "tmp")
	file, err := os.Create("tmp/config.yaml")
	if err != nil {
		t.Errorf("Failed to create tmp file %s", err.Error())
	}
	return file
}

func Cleanup(filename string) {
	os.Remove(filename)
}

func SetupCmd(command *cobra.Command) (*cobra.Command, *bytes.Buffer) {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(command)
	outputBuf := &bytes.Buffer{}
	rootCmd.SetOut(outputBuf)

	return rootCmd, outputBuf
}

func ExecuteCmd(t *testing.T, command *cobra.Command) {
	if err := command.Execute(); err != nil {
		t.Errorf("command failed: %v", err)
	}
}
