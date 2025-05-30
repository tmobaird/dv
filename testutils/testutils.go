package testutils

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func AssertEqual[T comparable](t *testing.T, expected, got T, msg ...interface{}) {
	t.Helper()
	if expected != got {
		t.Errorf("\nexpected '%v'\ngot '%v'\n%v", expected, got, msg)
	}
}

func AssertNotEqual[T comparable](t *testing.T, expected, got T, msg ...interface{}) {
	t.Helper()
	if expected == got {
		t.Errorf("expected '%v'\nto not equal '%v'\n%v", got, expected, msg)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("did not expect an error, got %s", err.Error())
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("expected an error")
	}
}

func FileContentEquals(t *testing.T, expected string, file *os.File) {
	t.Helper()

	_, err := file.Seek(0, 0)
	AssertNoError(t, err)
	data, err := io.ReadAll(file)
	AssertNoError(t, err)
	AssertEqual(t, expected, string(data))
}

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
	AssertNoError(t, err)
	return dirname
}

func RunCmd(t *testing.T, cmd *cobra.Command) *bytes.Buffer {
	outputBuf := &bytes.Buffer{}
	rootCmd, outputBuf := SetupCmd(cmd)
	rootCmd.SetArgs([]string{"list"})
	ExecuteCmd(t, rootCmd)
	return outputBuf
}

func CreateTodosFile(t *testing.T, filename, content string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	AssertNoError(t, err)
	_, err = file.Write([]byte(content))
	AssertNoError(t, err)
}
