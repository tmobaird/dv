package lg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tmobaird/dv/colors"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/lg/cmd"
	"github.com/tmobaird/dv/testutils"
)

func TestLg(t *testing.T) {
	t.Run("write creates a daily dev log for writing", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("EDITOR", "echo")
		rootCmd, outputBuf := testutils.SetupCmd(cmd.WriteCmd)
		rootCmd.SetArgs([]string{"write"})
		testutils.ExecuteCmd(t, rootCmd)

		filename := fmt.Sprintf("%s.md", time.Now().Format(core.LOG_FILE_TIME_FORMAT))
		testutils.AssertEqual(t, fmt.Sprintf("Opened file for writing %s.", filename), outputBuf.String())
		logFiles, err := os.ReadDir(core.LogDirectoryPath())
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, 1, len(logFiles))
		testutils.AssertEqual(t, filename, logFiles[0].Name())
	})

	t.Run("Show displays the latest dev log", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my daily dev log!")
		rootCmd, outputBuf := testutils.SetupCmd(cmd.ShowCmd)
		rootCmd.SetArgs([]string{"show"})
		testutils.ExecuteCmd(t, rootCmd)
		defer os.RemoveAll(core.LogDirectoryPath())

		testutils.AssertContains(t, outputBuf.String(), fmt.Sprintf("Date: %s", time.Now().Format(time.DateOnly)))
		testutils.AssertContains(t, outputBuf.String(), "latest")
		testutils.AssertContains(t, outputBuf.String(), "Here's my daily dev log!")
	})

	t.Run("Log displays an ordered list of dev logs", func(t *testing.T) {
		configFile := testutils.CreateConfigFile(t)
		defer testutils.Cleanup(configFile.Name())
		config := core.Config{Context: "main"}
		core.SaveConfig(configFile, config)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		logOneDate := time.Now().Add(time.Hour * -24)
		logTwoDate := time.Now()
		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(logOneDate)), "Here's dev log 1")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(logTwoDate)), "Here's dev log 2")
		rootCmd, outputBuf := testutils.SetupCmd(cmd.LogCmd)
		rootCmd.SetArgs([]string{"log"})
		testutils.ExecuteCmd(t, rootCmd)
		defer os.RemoveAll(core.LogDirectoryPath())

		testutils.AssertContains(t, outputBuf.String(), fmt.Sprintf("Date: %s %s(latest)%s", logTwoDate.Format(time.DateOnly), colors.EscapeCode(colors.CODE_BOLD), colors.EscapeCode(colors.CODE_RESET_BOLD)))
		testutils.AssertContains(t, outputBuf.String(), "Here's dev log 2")
		testutils.AssertContains(t, outputBuf.String(), fmt.Sprintf("Date: %s", logOneDate.Format(time.DateOnly)))
		testutils.AssertContains(t, outputBuf.String(), "Here's dev log 1")
	})
}
