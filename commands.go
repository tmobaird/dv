package main

import (
	"time"

	"github.com/google/uuid"
)

func add(args []string) ([]Todo, error) {
	// deserialize data
	todos, err := readJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	// add todo to list
	newTodo := Todo{
		Id:        uuid.New(),
		Name:      args[0],
		CreatedAt: time.Now().String(),
	}

	todos = append(todos, newTodo)
	writeErr := writeTodosToFile(todos)
	if writeErr != nil {
		todos = []Todo{}
		return todos, err
	}

	return todos, nil
}

func list() ([]Todo, error) {
	todos, err := readJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	return todos, nil
}

func delete(uid string) ([]Todo, error) {
	todos, err := readJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	var newTodos []Todo
	for _, todo := range todos {
		if todo.Id.String() != uid {
			newTodos = append(newTodos, todo)
		}
	}

	writeErr := writeTodosToFile(newTodos)

	if writeErr != nil {
		todos = []Todo{}
		return todos, err
	}

	return newTodos, nil
}
