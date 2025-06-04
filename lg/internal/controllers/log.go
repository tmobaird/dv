package controllers

import (
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/tmobaird/dv/core"
)

type LogArgs struct {
	Before string
	After  string
}

func (args LogArgs) BeforeDate() (time.Time, error) {
	return time.Parse(time.DateOnly, args.Before)
}

func (args LogArgs) AfterDate() (time.Time, error) {
	return time.Parse(time.DateOnly, args.After)
}

func (args LogArgs) HasBefore() bool {
	return args.Before != ""
}

func (args LogArgs) HasAfter() bool {
	return args.After != ""
}

func (controller Controller) RunLog(args LogArgs) (string, error) {
	before, after, err := parseBeforeAndAfter(args)
	if err != nil {
		return "", err
	}

	cmd, stdin, err := startPager()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(core.LogDirectoryPath())
	if err != nil {
		cmd.Cancel()
		stdin.Close()
		return "", err
	}

	filtered := filterLogs(entries, before, after)
	sortLogs(filtered)

	for i, entry := range filtered {
		output := logEntryOutput(entry.Name(), i == 0)
		stdin.Write([]byte(output))
	}

	stdin.Close()
	cmd.Wait()

	return "", nil
}

func parseBeforeAndAfter(args LogArgs) (time.Time, time.Time, error) {
	var err error
	before := time.Now().Add(24 * time.Hour)
	after := time.Now().Add(-1 * 20000 * 24 * time.Hour)
	if args.HasBefore() {
		before, err = args.BeforeDate()
		if err != nil {
			return before, after, err
		}
	}
	if args.HasAfter() {
		after, err = args.AfterDate()
		if err != nil {
			return before, after, err
		}
	}

	return before, after, nil
}

func getPager() string {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}
	return pager
}

func filterLogs(files []os.DirEntry, before time.Time, after time.Time) []os.DirEntry {
	filtered := []os.DirEntry{}
	for _, entry := range files {
		t, err := time.Parse(core.LOG_FILE_TIME_FORMAT, strings.Split(entry.Name(), ".md")[0])
		if err == nil && entry.Type().IsRegular() && t.Before(before) && t.After(after) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func sortLogs(files []os.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		a, _ := time.Parse(core.LOG_FILE_TIME_FORMAT, strings.Split(files[i].Name(), ".md")[0])
		b, _ := time.Parse(core.LOG_FILE_TIME_FORMAT, strings.Split(files[j].Name(), ".md")[0])

		return a.After(b)
	})
}

func startPager() (*exec.Cmd, io.WriteCloser, error) {
	pager := getPager()

	cmd := exec.Command(pager)
	cmd.Stdout = os.Stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return cmd, stdin, err
	}

	if err := cmd.Start(); err != nil {
		return cmd, stdin, err
	}

	return cmd, stdin, nil
}
