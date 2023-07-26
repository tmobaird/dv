package main

import (
	"errors"
	"strconv"
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

func moveTodoToPosition(todos []Todo, uid string, position int) []Todo {
	var currentPosition int
	for i, obj := range todos {
        if obj.Id.String() == uid || strconv.Itoa(i + 1) == uid {
            currentPosition = i
            break
        }
    }

	if currentPosition == position || currentPosition >= len(todos) {
        return todos
    }

	var newTodos []Todo
	movedObject := todos[currentPosition]
    newTodos = append(todos[:currentPosition], todos[currentPosition+1:]...)

	if position >= len(todos) {
        // If the newPosition is greater than the length of the array, append the object
        newTodos = append(todos, movedObject)
    } else {
        // Otherwise, insert the object at the newPosition
        newTodos = append(newTodos[:position], append([]Todo{movedObject}, newTodos[position:]...)...)
    }

    return newTodos
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
	for index, todo := range todos {
		todoIndex := strconv.Itoa(index + 1)
		if !isInArray(todo.Id.String(), uids) && !isInArray(todoIndex, uids) {
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
	for index, todo := range todos {
		todoIndex := strconv.Itoa(index + 1)
		if todo.Id.String() == uid || todoIndex == uid {
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
	for index, todo := range todos {
		todoIndex := strconv.Itoa(index + 1)
		if todo.Id.String() == uid || todoIndex == uid {
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
	for index, todo := range todos {
		todoIndex := strconv.Itoa(index + 1)
		if todo.Id.String() == uid || todoIndex == uid {
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

func rank(uid, rank string, r ReaderWriter) ([]Todo, error) {
	todos, readErr := r.ReadJSONFileToMap()

	if readErr != nil {
		todos = []Todo{}
		return todos, readErr
	}

	rankInt, rankErr := strconv.Atoi(rank)

	if rankErr != nil {
		todos = []Todo{}
		return todos, rankErr
	}

	if rankInt > len(todos) || rankInt < 1 {
		todos = []Todo{}
		return todos, errors.New("Position is out of range")
	}

	newTodos := moveTodoToPosition(todos, uid, rankInt - 1)

	writeErr := r.WriteTodosToFile(newTodos)
	if writeErr != nil {
		todos = []Todo{}
		return todos, writeErr
	}

	return newTodos, nil
}