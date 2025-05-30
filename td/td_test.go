package td

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/cmd"
	"github.com/tmobaird/dv/testutils"
)

func TestTdIntegration(t *testing.T) {
	t.Run("td config", func(t *testing.T) {
		file := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(file.Name())

		config := core.Config{Context: "main"}
		core.SaveConfig(file, config)

		rootCmd, outputBuf := testutils.SetupCmd(cmd.ConfigCmd)
		rootCmd.SetArgs([]string{"config"})
		testutils.ExecuteCmd(t, rootCmd)

		got := outputBuf.String()
		if !strings.Contains(got, "Context: main") {
			t.Errorf("Expected %s to contain %s", got, "Context: main")
		}
	})

	t.Run("td context", func(t *testing.T) {
		file := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(file.Name())

		config := core.Config{Context: "my-context"}
		core.SaveConfig(file, config)

		rootCmd, outputBuf := testutils.SetupCmd(cmd.ContextCmd)
		rootCmd.SetArgs([]string{"context"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "my-context"
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td context new-value", func(t *testing.T) {
		file := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(file.Name())

		config := core.Config{Context: "my-context"}
		core.SaveConfig(file, config)

		rootCmd, outputBuf := testutils.SetupCmd(cmd.ContextCmd)
		rootCmd.SetArgs([]string{"context", "updated-context"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "Updated context to updated-context"
		got := outputBuf.String()

		testutils.AssertEqual(t, expected, got)
		loaded, err := core.ReadConfig(os.DirFS(core.BasePath()))
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, loaded.Context, "updated-context")
	})

	t.Run("td list", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)
		testutils.CreateTodosFile(t, "tmp/lists/main.md", "- [ ] todo numero uno")

		rootCmd, outputBuf := testutils.SetupCmd(cmd.ListCmd)
		rootCmd.SetArgs([]string{"list"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "1. [ ] todo numero uno\n"
		got := outputBuf.String()

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td open", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		rootCmd, outputBuf := testutils.SetupCmd(cmd.OpenCmd)
		rootCmd.SetArgs([]string{"open"})
		testutils.ExecuteCmd(t, rootCmd)

		testutils.AssertEqual(t, "Must set $EDITOR to edit config", outputBuf.String())
	})

	t.Run("td add", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)

		rootCmd, outputBuf := testutils.SetupCmd(cmd.AddCmd)
		rootCmd.SetArgs([]string{"add", "do homework"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "\"do homework\" added to list."
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)

		output := testutils.RunCmd(t, cmd.ListCmd)
		expected = "1. [ ] do homework\n"
		got = output.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td rm", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)
		testutils.CreateTodosFile(t, "tmp/lists/main.md", "- [ ] Todo One")

		rootCmd, outputBuf := testutils.SetupCmd(cmd.RemoveCmd)
		rootCmd.SetArgs([]string{"rm", "1"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "\"Todo One\" removed from list."
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)

		output := testutils.RunCmd(t, cmd.ListCmd)
		expected = "No todos in list."
		got = output.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td mv", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)
		testutils.CreateTodosFile(t, "tmp/lists/main.md", "- [ ] Todo One")

		rootCmd, outputBuf := testutils.SetupCmd(cmd.RenameCmd)
		rootCmd.SetArgs([]string{"mv", "1", "updated todo name"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "\"Todo One\" updated to \"updated todo name\"."
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)

		output := testutils.RunCmd(t, cmd.ListCmd)
		expected = "1. [ ] updated todo name\n"
		got = output.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td rank", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)
		testutils.CreateTodosFile(t, "tmp/lists/main.md", "- [ ] Todo A\n- [ ] Todo B\n")

		rootCmd, outputBuf := testutils.SetupCmd(cmd.RankCmd)
		rootCmd.SetArgs([]string{"rank", "1", "2"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "\"Todo A\" ranked to index 2."
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)

		output := testutils.RunCmd(t, cmd.ListCmd)
		expected = "1. [ ] Todo B\n2. [ ] Todo A\n"
		got = output.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td mark", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)
		testutils.CreateTodosFile(t, "tmp/lists/main.md", "- [ ] Todo A\n- [ ] Todo B\n")

		rootCmd, outputBuf := testutils.SetupCmd(cmd.MarkCmd)
		rootCmd.SetArgs([]string{"mark", "1", "done"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := "\"Todo A\" marked done."
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)

		output := testutils.RunCmd(t, cmd.ListCmd)
		expected = "1. [x] Todo A\n2. [ ] Todo B\n"
		got = output.String()
		testutils.AssertEqual(t, expected, got)
	})

	t.Run("td schedule", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateListsDirectory(t)
		defer os.RemoveAll(dirname)
		defer os.RemoveAll("tmp/schedules")
		testutils.CreateTodosFile(t, "tmp/lists/main.md", "- [ ] Todo A\n- [ ] Todo B\n")

		rootCmd, outputBuf := testutils.SetupCmd(cmd.ScheduleCmd)
		rootCmd.SetArgs([]string{"schedule", "--no-calendar"})
		testutils.ExecuteCmd(t, rootCmd)

		expected := fmt.Sprintf(
			"Here is your schedule for %s\n%s-%s: Todo A\n%s-%s: Todo B\n",
			time.Now().Format(time.DateOnly),
			time.Now().Format(time.Kitchen),
			time.Now().Add(30*time.Minute).Format(time.Kitchen),
			time.Now().Add(30*time.Minute).Format(time.Kitchen),
			time.Now().Add(60*time.Minute).Format(time.Kitchen),
		)
		got := outputBuf.String()
		testutils.AssertEqual(t, expected, got)
	})
}
