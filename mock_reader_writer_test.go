package main

import (
	"errors"
	"os"
)

type MockReaderWriter struct {
	todos []Todo
	MkdirAllFunc func(path string, perm os.FileMode) error
	StatFunc func(name string) (os.FileInfo, error)
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
}

func (m *MockReaderWriter) ReadJSONFileToMap() ([]Todo, error) {
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

func (m *MockReaderWriter) ReadConfigFile() (Config, error) {
	return Config{}, nil
}

type ErrorMockReader struct{
	todos []Todo
}

func (m *ErrorMockReader) ReadJSONFileToMap() ([]Todo, error) {
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

func (m *ErrorMockReader) ReadConfigFile() (Config, error) {
	return Config{}, nil
}

type ErrorMockWriter struct{
	todos []Todo
}

func (m *ErrorMockWriter) ReadJSONFileToMap() ([]Todo, error) {
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

func (m *ErrorMockWriter) ReadConfigFile() (Config, error) {
	return Config{}, nil
}
