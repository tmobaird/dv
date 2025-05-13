package models

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"td/internal"
)

type TodoType int

const (
	Blank TodoType = iota
	Deep
	Shallow
	Quick
)

func GetTodoType(num TodoType) string {
	switch num {
	case 0:
		return "Blank"
	case 1:
		return "Deep"
	case 2:
		return "Shallow"
	case 3:
		return "Quick"
	default:
		return ""
	}
}

const MINUTE_UNIT = "minute"
const HOUR_UNIT = "hour"
const BLOCK_UNIT = "block"

type Metadata struct {
	DurationValue int
	DurationUnit  string
	Type          TodoType
}

type Todo struct {
	Name     string
	Complete bool
	Metadata Metadata
}

func (todo *Todo) ToMd() string {
	completeChar := " "
	if todo.Complete {
		completeChar = "x"
	}
	return fmt.Sprintf("- [%s] %s\n", completeChar, todo.Name)
}

func (todo *Todo) Status() string {
	if todo.Complete {
		return "done"
	} else {
		return "not done"
	}
}

func GetAllTodos(filename string) ([]Todo, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return []Todo{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return []Todo{}, err
	}

	todos := parseTodos(data)
	return todos, nil
}

func WriteTodos(context string, todos []Todo) error {
	file, err := os.OpenFile(internal.TodoFilePath(context), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content := ""
	for _, todo := range todos {
		content += todo.ToMd()
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}

func DefaultMetadata() Metadata {
	return Metadata{DurationValue: 1, DurationUnit: BLOCK_UNIT, Type: Blank}
}

func parseTodos(bytes []byte) []Todo {
	todos := []Todo{}
	data := string(bytes)
	initialRegex := regexp.MustCompile("^- .*")                                   // checks to make sure it is a valid list item
	partsRegex := regexp.MustCompile(`^- \[([^\]]+)\] (.*?)\s*(?:\(([^)]*)\))?$`) // captures complete and name
	for _, line := range strings.Split(data, "\n") {
		if initialRegex.Match([]byte(line)) {
			parts := partsRegex.FindStringSubmatch(line)
			complete := false
			if len(parts) >= 3 {
				if parts[1] == "x" {
					complete = true
				}

				name := strings.TrimSpace(parts[2])
				metadata := DefaultMetadata()
				if len(parts) >= 4 {
					metadata = parseMetadata(parts[3])
				}
				todos = append(todos, Todo{Name: name, Complete: complete, Metadata: metadata})
			}
		}
	}
	return todos
}

func FilterTodos(todos []Todo, hideCompleted bool) []Todo {
	result := []Todo{}
	for _, todo := range todos {
		if hideCompleted && todo.Complete {
			// do nothing, ie filter out
		} else {
			result = append(result, todo)
		}
	}
	return result
}

func parseMetadata(metadataString string) Metadata {
	metadata := DefaultMetadata()
	attrs := strings.Split(metadataString, ",")
	for _, attr := range attrs {
		if strings.Contains(attr, "=") {
			parts := strings.Split(attr, "=")
			key := parts[0]
			value := parts[1]

			if key == "duration" {
				value, unit, err := parseDuration(value)
				if err == nil {
					metadata.DurationValue = value
					metadata.DurationUnit = unit
				}
			} else if key == "type" {
				metadata.Type = parseType(value)
			}
		}
	}
	return metadata
}

func parseDuration(durationString string) (int, string, error) {
	regex := regexp.MustCompile(`^([0-9]+)([h,m,b]{1})`)
	parts := regex.FindStringSubmatch(durationString)
	if len(parts) < 3 {
		return -1, "", fmt.Errorf("expected 2 parts, found %d", len(parts))
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, "", fmt.Errorf("expected a valid number, got %s", err.Error())
	}

	var units string
	switch parts[2] {
	case "h":
		units = HOUR_UNIT
	case "b":
		units = BLOCK_UNIT
	case "m":
		units = MINUTE_UNIT
	default:
		return -1, "", fmt.Errorf("expected units of one of the following: [h,b,m], got %s", parts[1])
	}

	return num, units, nil
}

func parseType(typeString string) TodoType {
	if typeString == "deep" {
		return Deep
	} else if typeString == "shallow" {
		return Shallow
	} else if typeString == "quick" {
		return Quick
	}
	return Blank
}
