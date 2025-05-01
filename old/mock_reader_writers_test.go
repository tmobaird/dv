package main

import (
	"errors"
	"os"
)

type MockReaderWriter struct {
	todos         []Todo
	Config        Config
	MkdirAllFunc  func(path string, perm os.FileMode) error
	StatFunc      func(name string) (os.FileInfo, error)
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
}

func (m *MockReaderWriter) ReadTodosFromFile() ([]Todo, error) {
	return m.todos, nil
}

func (m *MockReaderWriter) WriteTodosToFile(todos []Todo) error {
	m.todos = todos
	return nil
}

func (m *MockReaderWriter) EnsureTodosFileExists() error {
	return nil
}

func (m *MockReaderWriter) EnsureConfigFileExists() error {
	return nil
}

func (m *MockReaderWriter) ReadConfigFromFile() (Config, error) {
	return m.Config, nil
}

func (m *MockReaderWriter) WriteConfigToFile(config Config) error {
	m.Config = config
	return nil
}

type ErrorMockReader struct {
	todos []Todo
}

func (m *ErrorMockReader) ReadTodosFromFile() ([]Todo, error) {
	return nil, errors.New("Failed to read file")
}

func (m *ErrorMockReader) WriteTodosToFile(todos []Todo) error {
	return nil
}

func (m *ErrorMockReader) EnsureTodosFileExists() error {
	return nil
}

func (m *ErrorMockReader) EnsureConfigFileExists() error {
	return nil
}

func (m *ErrorMockReader) ReadConfigFromFile() (Config, error) {
	return Config{}, errors.New("Failed to read file")
}

func (m *ErrorMockReader) WriteConfigToFile(config Config) error {
	return nil
}

type ErrorMockWriter struct {
	todos []Todo
}

func (m *ErrorMockWriter) ReadTodosFromFile() ([]Todo, error) {
	return m.todos, nil
}

func (m *ErrorMockWriter) WriteTodosToFile(todos []Todo) error {
	return errors.New("Failed to write file")
}

func (m *ErrorMockWriter) EnsureTodosFileExists() error {
	return nil
}

func (m *ErrorMockWriter) EnsureConfigFileExists() error {
	return nil
}

func (m *ErrorMockWriter) ReadConfigFromFile() (Config, error) {
	return Config{}, nil
}

func (m *ErrorMockWriter) WriteConfigToFile(config Config) error {
	return errors.New("Failed to write file")
}