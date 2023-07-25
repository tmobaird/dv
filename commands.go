package main

import (
	"time"

	"github.com/google/uuid"
)

func isInArray(target string, arr []string) bool {
    for _, item := range arr {
        if item == target {
            return true
        }
    }
    return false
}


func add(args []string, r ReaderWriter) ([]Todo, error) {
	// deserialize data
	todos, err := r.ReadJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	for _, arg := range args {
		newTodo := Todo{
			Id:        uuid.New(),
			Name:      arg,
			CreatedAt: time.Now().String(),
		}
		todos = append(todos, newTodo)
	}

	writeErr := r.WriteTodosToFile(todos)
	if writeErr != nil {
		todos = []Todo{}
		return todos, writeErr
	}

	return todos, nil
}

func list(r ReaderWriter) ([]Todo, error) {
	todos, err := r.ReadJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	return todos, nil
}

func delete(uids []string, r ReaderWriter) ([]Todo, error) {
	todos, err := r.ReadJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	var newTodos []Todo
	for _, todo := range todos {
		if !isInArray(todo.Id.String(), uids) {
			newTodos = append(newTodos, todo)
		}
	}

	writeErr := r.WriteTodosToFile(newTodos)

	if writeErr != nil {
		todos = []Todo{}
		return todos, writeErr
	}

	return newTodos, nil
}

func done(uid string, r ReaderWriter) ([]Todo, error) {
	todos, err := r.ReadJSONFileToMap()
	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	var newTodos []Todo
	for _, todo := range todos {
		if todo.Id.String() == uid {
			todo.Done = true
		}
		newTodos = append(newTodos, todo)
	}

	writeErr := r.WriteTodosToFile(newTodos)
	if writeErr != nil {
		todos = []Todo{}
		return todos, writeErr
	}

	return newTodos, nil
}

func undo(uid string, r ReaderWriter) ([]Todo, error) {
	todos, err := r.ReadJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	var newTodos []Todo
	for _, todo := range todos {
		if todo.Id.String() == uid {
			todo.Done = false
		}
		newTodos = append(newTodos, todo)
	}

	writeErr := r.WriteTodosToFile(newTodos)

	if writeErr != nil {
		todos = []Todo{}
		return todos, writeErr
	}

	return newTodos, nil
}

func edit(uid, name string, r ReaderWriter) ([]Todo, error) {
	todos, err := r.ReadJSONFileToMap()

	if err != nil {
		todos = []Todo{}
		return todos, err
	}

	var newTodos []Todo
	for _, todo := range todos {
		if todo.Id.String() == uid {
			todo.Name = name
		}
		newTodos = append(newTodos, todo)
	}

	writeErr := r.WriteTodosToFile(newTodos)
	if writeErr != nil {
		todos = []Todo{}
		return todos, writeErr
	}

	return newTodos, nil
}