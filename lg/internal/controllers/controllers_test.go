package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/testutils"
)

func TestControllers(t *testing.T) {
	t.Run("#RunShow returns latest dev log", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my daily dev log!")

		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunShow()

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my daily dev log!")
	})

	t.Run("#RunShow returns error when latest dev log does not exist", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")

		_, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunShow()

		testutils.AssertError(t, err)
		testutils.AssertContains(t, err.Error(), "latest dev log does not exist")
	})

	t.Run("#RunShow returns today's dev log when given 'today'", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my daily dev log!")
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{"today"}}.RunShow()

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my daily dev log!")
	})

	t.Run("#RunShow return error when given 'today' and no dev log exists", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")

		_, err := Controller{Config: core.Config{Context: "main"}, Args: []string{"today"}}.RunShow()

		testutils.AssertError(t, err)
		testutils.AssertContains(t, err.Error(), fmt.Sprintf("%s.md does not exist", time.Now().Format(core.LOG_FILE_TIME_FORMAT)))
	})

	t.Run("#RunShow returns yesterday's dev log when given 'yesterday'", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-24*time.Hour))), "Here's my daily dev log!")
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{"yesterday"}}.RunShow()

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my daily dev log!")
	})

	t.Run("#RunShow return error when given 'yesterday' and no dev log exists", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my daily dev log!")
		_, err := Controller{Config: core.Config{Context: "main"}, Args: []string{"yesterday"}}.RunShow()

		testutils.AssertError(t, err)
		testutils.AssertContains(t, err.Error(), fmt.Sprintf("%s.md does not exist", time.Now().Add(-24*time.Hour).Format(core.LOG_FILE_TIME_FORMAT)))
	})

	t.Run("#RunShow returns dev log from 2 days ago when given '2d'", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-48*time.Hour))), "Here's my daily dev log!")
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{"2d"}}.RunShow()

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my daily dev log!")
	})

	t.Run("#RunShow return error when given '2d' and no dev log exists", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-24*time.Hour))), "Here's my daily dev log!")
		_, err := Controller{Config: core.Config{Context: "main"}, Args: []string{"2d"}}.RunShow()

		testutils.AssertError(t, err)
		testutils.AssertContains(t, err.Error(), fmt.Sprintf("%s.md does not exist", time.Now().Add(-48*time.Hour).Format(core.LOG_FILE_TIME_FORMAT)))
	})

	t.Run("#RunLog returns dev logs in the correct order", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-24*time.Hour))), "Here's my first daily dev log!")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my second daily dev log!")
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunLog(LogArgs{})

		testutils.AssertNoError(t, err)
		keyLines := strings.Split(output.String(), "\n\n")
		testutils.AssertContains(t, keyLines[0], time.Now().Format(time.DateOnly))
		testutils.AssertContains(t, keyLines[1], "Here's my second daily dev log!")
		testutils.AssertContains(t, keyLines[2], time.Now().Add(-24*time.Hour).Format(time.DateOnly))
		testutils.AssertContains(t, keyLines[3], "Here's my first daily dev log!")
	})

	t.Run("#RunLog can handle before filter", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-24*time.Hour))), "Here's my first daily dev log!")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my second daily dev log!")
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunLog(LogArgs{Before: time.Now().Format(time.DateOnly)})

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my first daily dev log!")
		testutils.AssertDoesNotContain(t, output.String(), "Here's my first dev daily dev log!")
	})

	t.Run("#RunLog can handle after filter", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-24*time.Hour))), "Here's my first daily dev log!")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my second daily dev log!")
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunLog(LogArgs{After: time.Now().Add(-24 * time.Hour).Format(time.DateOnly)})

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my second daily dev log!")
		testutils.AssertDoesNotContain(t, output.String(), "Here's my first daily dev log!")
	})

	t.Run("#RunLog can handle before and after filters at the same time", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-48*time.Hour))), "Here's my first daily dev log!")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now().Add(-24*time.Hour))), "Here's my second daily dev log!")
		testutils.CreateLogFile(t, filepath.Join("tmp/logs", core.LogFileName(time.Now())), "Here's my third daily dev log!")
		logArgs := LogArgs{
			Before: time.Now().Format(time.DateOnly),
			After:  time.Now().Add(-48 * time.Hour).Format(time.DateOnly),
		}
		output, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunLog(logArgs)

		testutils.AssertNoError(t, err)
		testutils.AssertContains(t, output.String(), "Here's my second daily dev log!")
		testutils.AssertDoesNotContain(t, output.String(), "Here's my first daily dev log!")
		testutils.AssertDoesNotContain(t, output.String(), "Here's my third daily dev log!")
	})

	t.Run("#RunLog returns the correct thing when no dev logs exist", func(t *testing.T) {
		testutils.SetupDv(t)

		dirname := testutils.CreateLogsDirectory(t)
		defer os.RemoveAll(dirname)

		os.Setenv("PAGER", "this-is-a-fake-tool")
		_, err := Controller{Config: core.Config{Context: "main"}, Args: []string{}}.RunLog(LogArgs{})

		testutils.AssertError(t, err)
		testutils.AssertEqual(t, "no dev logs exist", err.Error())
	})
}
