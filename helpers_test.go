package main

import (
	"bytes"
	"os"
	"td/cmd"
	"td/internal/testutils"
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

func CreateListsDirectory(t *testing.T) string {
	dirname := "tmp/lists"
	err := os.MkdirAll(dirname, 0755)
	testutils.AssertNoError(t, err)
	return dirname
}

func RunListCmd(t *testing.T) *bytes.Buffer {
	outputBuf := &bytes.Buffer{}
	rootCmd, outputBuf := SetupCmd(cmd.ListCmd)
	rootCmd.SetArgs([]string{"list"})
	ExecuteCmd(t, rootCmd)
	return outputBuf
}

func CreateTodosFile(t *testing.T, filename, content string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	testutils.AssertNoError(t, err)
	_, err = file.Write([]byte(content))
	testutils.AssertNoError(t, err)
}
