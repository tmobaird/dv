package controllers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/tmobaird/dv/colors"
	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/lg/internal/markdown"
)

const LOG_FILE_TEMPLATE = `## Worked on


## Up next


## Issues & Surprises


## Other thoughts
`

type Controller struct {
	Args   []string
	Config core.Config
}

func (controller Controller) RunWrite() (OutputTarget, error) {
	writeTo := NewSimpleOutputTarget()
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return writeTo, errors.New("must set $EDITOR to edit config")
	}

	logFileName := core.LogFileName(time.Now())
	logFilePath, err := createLogFile(time.Now())
	if err != nil {
		return writeTo, err
	}

	cmd := exec.Command(os.Getenv("EDITOR"), logFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return writeTo, err
	}

	writeTo.WriteString(fmt.Sprintf("Opened file for writing %s.", logFileName))
	return writeTo, nil
}

func (controller Controller) RunShow() (OutputTarget, error) {
	writeTo := OutputTarget{}
	var err error
	on := time.Now()
	if len(controller.Args) > 0 {
		term := controller.Args[0]
		if term == "today" {
			on = time.Now()
		} else if term == "yesterday" {
			on = time.Now().Add(time.Hour * 24 * -1)
		} else {
			regex := regexp.MustCompile(`^(\d+)d`)
			parts := regex.FindStringSubmatch(term)
			if len(parts) < 2 {
				fmt.Printf("Warning: Bad indicator \"%s\". Using today\n", term)
			} else {
				daysAgo, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Printf("Warning: Bad indicator \"%s\". Using today\n", term)
				} else {
					on = time.Now().Add(time.Hour * (24 * time.Duration(daysAgo)) * -1)
				}
			}
		}
	} else { // use latest
		on, err = getLatestLogDate()
		if err != nil {
			return writeTo, err
		}
	}

	cmd, stdin, err := startPager()
	if err != nil {
		writeTo = NewSimpleOutputTarget()
	} else {
		writeTo = NewCmdOutputTarget(stdin, cmd)
	}

	if core.LogFileExists(on) {
		contents := logEntryOutput(core.LogFileName(on), true)
		writeTo.WriteString(contents)
		writeTo.Close()
		err = nil
	} else {
		err = fmt.Errorf("%s does not exist", core.LogFileName(on))
	}

	return writeTo, err
}

func createLogFile(t time.Time) (string, error) {
	err := createLogsDirectory()
	if err != nil {
		return "", err
	}

	logFilePath := core.LogFilePath(t)
	if !core.LogFileExists(t) {
		f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}

		_, err = f.Write([]byte(LOG_FILE_TEMPLATE))
		if err != nil {
			return "", err
		}

		defer f.Close()
	}

	return logFilePath, nil
}

func createLogsDirectory() error {
	dirPath := core.LogDirectoryPath()
	return os.MkdirAll(dirPath, 0755)
}

func getLatestLogDate() (time.Time, error) {
	entries, err := os.ReadDir(core.LogDirectoryPath())
	if err != nil {
		return time.Now(), err
	}

	sort.Slice(entries, func(i, j int) bool {
		a, _ := core.LogFileNameToTime(entries[i].Name())
		b, _ := core.LogFileNameToTime(entries[j].Name())

		return a.After(b)
	})

	if len(entries) == 0 {
		return time.Now(), errors.New("latest dev log does not exist")
	}

	on, err := core.LogFileNameToTime(entries[0].Name())
	if err != nil {
		return time.Now(), err
	}

	return on, nil
}

func logEntryOutput(filename string, latest bool) string {
	day, _ := core.LogFileNameToTime(filename)
	contents, err := os.ReadFile(filepath.Join(core.LogDirectoryPath(), filename))
	if err == nil {
		prefix := ""
		if latest {
			prefix += " " + colors.AddTextStyle("(latest)", colors.CODE_BOLD)
		}
		prefix = colors.AddColor(fmt.Sprintf("Date: %s%s", day.Format(time.DateOnly), prefix), colors.FG_YELLOW)
		return fmt.Sprintf("%s\n\n%s\n", prefix, markdown.MdToOutput(string(contents)))
	}
	return ""
}
