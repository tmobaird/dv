package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func readJSONFileToMap() ([]Todo, error) {
	raw, err := os.ReadFile("./tmp/todos.json")
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

func writeTodosToFile(todos []Todo) error {
	// serialize data into json
	jsonData, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		return err
	}

	// write to json file
	if err := ioutil.WriteFile("./tmp/todos.json", jsonData, 0644); err != nil {
		return err
	}

	return nil
}