package main

import (
	"encoding/json"
	"os"
)

var storageDirectory string = ".td"
var storageFileName string = "todos.json"
var configFileName string = "config.json"

type ReaderWriter interface {
	ReadTodosFromFile() ([]Todo, error)
	WriteTodosToFile([]Todo) error
	EnsureTodosFileExists() error
	EnsureConfigFileExists() error
	ReadConfigFromFile() (Config, error)
	WriteConfigToFile(Config) error
}

type TdReaderWriter struct {
	Context       string
	MkdirAllFunc  func(path string, perm os.FileMode) error
	StatFunc      func(name string) (os.FileInfo, error)
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
	ReadFileFunc  func(filename string) ([]byte, error)
}

func (r TdReaderWriter) todosFilePath() string {
	return r.Context + "/" + storageFileName
}

func (r TdReaderWriter) configFilePath() string {
	return r.Context + "/" + configFileName
}

func (r *TdReaderWriter) EnsureTodosFileExists() error {
	if err := r.MkdirAllFunc(r.Context, 0755); err != nil {
		return err
	}

	if _, err := r.StatFunc(r.todosFilePath()); os.IsNotExist(err) {
		r.WriteTodosToFile([]Todo{})
	}
	return nil
}

func (r *TdReaderWriter) EnsureConfigFileExists() error {
	if err := r.MkdirAllFunc(r.Context, 0755); err != nil {
		return err
	}

	if _, err := r.StatFunc(r.configFilePath()); os.IsNotExist(err) {
		r.WriteFileFunc(r.configFilePath(), []byte("{}"), 0644)
	}
	return nil
}

func (r *TdReaderWriter) ReadTodosFromFile() ([]Todo, error) {
	raw, err := r.ReadFileFunc(r.todosFilePath())
	if err != nil {
		return nil, err
	}
	var data []Todo
	marshalErr := json.Unmarshal(raw, &data)

	if marshalErr != nil {
		return nil, marshalErr
	}

	return data, nil
}

func (r *TdReaderWriter) WriteTodosToFile(todos []Todo) error {
	// serialize data into json
	jsonData, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		return err
	}

	// write to json file
	if err := r.WriteFileFunc(r.todosFilePath(), jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (r *TdReaderWriter) ReadConfigFromFile() (Config, error) {
	raw, err := r.ReadFileFunc(r.configFilePath())
	if err != nil {
		return Config{}, err
	}
	var config Config
	marshalErr := json.Unmarshal(raw, &config)

	if marshalErr != nil {
		return Config{}, marshalErr
	}

	return config, nil
}

func (r *TdReaderWriter) WriteConfigToFile(config Config) error {
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	err = r.WriteFileFunc(r.configFilePath(), jsonData, 0644)
	return err
}
