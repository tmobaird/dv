package main

import (
	"os"
	"strings"
	"td/cmd"
	"td/internal"
	"td/internal/testutils"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Run("td config", func(t *testing.T) {
		file := CreateConfigFile(t)
		defer Cleanup(file.Name())

		config := internal.Config{Context: "main"}
		internal.Save(file, config)

		rootCmd, outputBuf := SetupCmd(cmd.ConfigCmd)
		rootCmd.SetArgs([]string{"config"})
		ExecuteCmd(t, rootCmd)

		got := outputBuf.String()
		if !strings.Contains(got, "Context: main") {
			t.Errorf("Expected %s to contain %s", got, "Context: main")
		}
	})

	t.Run("td context", func(t *testing.T) {
		file := CreateConfigFile(t)
		defer Cleanup(file.Name())

		config := internal.Config{Context: "my-context"}
		internal.Save(file, config)

		rootCmd, outputBuf := SetupCmd(cmd.ContextCmd)
		rootCmd.SetArgs([]string{"context"})
		ExecuteCmd(t, rootCmd)

		expected := "my-context"
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td context new-value", func(t *testing.T) {
		file := CreateConfigFile(t)
		defer Cleanup(file.Name())

		config := internal.Config{Context: "my-context"}
		internal.Save(file, config)

		rootCmd, outputBuf := SetupCmd(cmd.ContextCmd)
		rootCmd.SetArgs([]string{"context", "updated-context"})
		ExecuteCmd(t, rootCmd)

		expected := "Updated context to updated-context"
		got := outputBuf.String()

		testutils.AssertEqual(t, expected, got)
		loaded, err := internal.Read(os.DirFS(internal.BasePath()))
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, loaded.Context, "updated-context")
	})

	t.Run("td list", func(t *testing.T) {
		configFile := CreateConfigFile(t)
		defer Cleanup(configFile.Name())
		config := internal.Config{Context: "main"}
		internal.Save(configFile, config)

		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		todosFile, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		todosFile.Write([]byte("- [ ] todo numero uno"))
		defer os.RemoveAll(dirname)

		rootCmd, outputBuf := SetupCmd(cmd.ListCmd)
		rootCmd.SetArgs([]string{"list"})
		ExecuteCmd(t, rootCmd)

		expected := "1. [ ] todo numero uno\n"
		got := outputBuf.String()

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td open", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		rootCmd, outputBuf := SetupCmd(cmd.OpenCmd)
		rootCmd.SetArgs([]string{"open"})
		ExecuteCmd(t, rootCmd)

		testutils.AssertEqual(t, "Must set $EDITOR to edit config", outputBuf.String())
	})

	t.Run("td add", func(t *testing.T) {
		configFile := CreateConfigFile(t)
		defer Cleanup(configFile.Name())
		config := internal.Config{Context: "main"}
		internal.Save(configFile, config)

		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		defer os.RemoveAll(dirname)

		rootCmd, outputBuf := SetupCmd(cmd.AddCmd)
		rootCmd.SetArgs([]string{"add", "do homework"})
		ExecuteCmd(t, rootCmd)

		expected := "\"do homework\" added to list."
		got := outputBuf.String()

		testutils.AssertEqual(t, expected, got)

		rootCmd, outputBuf = SetupCmd(cmd.ListCmd)
		rootCmd.SetArgs([]string{"list"})
		ExecuteCmd(t, rootCmd)

		expected = "1. [ ] do homework\n"
		got = outputBuf.String()

		testutils.AssertEqual(t, expected, got)
	})
}
