package main

import (
	"errors"
	"os"
	"testing"
)

func TestEnsureTodosFileExists(t *testing.T) {
	t.Run("creates the .td directory and todos file if it doesn't exist", func(t *testing.T) {
		tdCreated := false
		todosCreated := false
		r := &RealReaderWriter{
			MkdirAllFunc: func(path string, perm os.FileMode) error {
				tdCreated = true
				return nil
			},
			StatFunc: func(name string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
				todosCreated = true
				return nil
			},
		}
		err := r.EnsureTodosFileExists()

		assertNoError(t, err)

		assertTdCreated(t, tdCreated, true)
		assertTodosFileCreated(t, todosCreated, true)
	})

	t.Run("does not create the todos file if it already exists", func(t *testing.T) {
		todosCreated := false
		r := &RealReaderWriter{
			MkdirAllFunc: func(path string, perm os.FileMode) error {
				return nil
			},
			StatFunc: func(name string) (os.FileInfo, error) {
				return &MockFileInfo{}, os.ErrNotExist
			},
			WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
				todosCreated = true
				return nil
			},
		}

		err := r.EnsureTodosFileExists()

		assertNoError(t, err)
		assertTodosFileCreated(t, todosCreated, true)
	})

	t.Run("returns error when unable to create .td directory", func(t *testing.T) {
		r := &RealReaderWriter{
			MkdirAllFunc: func(path string, perm os.FileMode) error {
				return errors.New("unable to create .td directory")
			},
		}

		err := r.EnsureTodosFileExists()

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestReadJSONFileToMap(t *testing.T) {
	t.Run("returns todos from todos.json", func(t *testing.T) {
		r := &RealReaderWriter{
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte(`[]`), nil
			},
		}

		todos, err := r.ReadJSONFileToMap()

		assertNoError(t, err)
		if len(todos) != 0 {
			t.Errorf("Expected todos to be empty, got %v", todos)
		}
	})

	t.Run("returns error when read fails", func(t *testing.T) {
		var want error = errors.New("unable to read file")
		r := &RealReaderWriter{
			ReadFileFunc: func(filename string) ([]byte, error) {
				return nil, want
			},
		}

		_, got := r.ReadJSONFileToMap()

		assertError(t, got, want)
	})

	t.Run("returns error when unmarshalling fails", func(t *testing.T) {
		r := &RealReaderWriter{
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte(`[ded]`), nil
			},
		}

		_, got := r.ReadJSONFileToMap()

		if got == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestWriteMapToJSONFile(t *testing.T) {
	t.Run("returns no error when write is successful", func(t *testing.T) {
		r := &RealReaderWriter{
			WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
				return nil
			},
		}

		err := r.WriteTodosToFile([]Todo{})

		assertNoError(t, err)
	})

	t.Run("returns error when write fails", func(t *testing.T) {
		var want error = errors.New("unable to write file")
		r := &RealReaderWriter{
			WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
				return want
			},
		}

		got := r.WriteTodosToFile([]Todo{})

		assertError(t, got, want)
	})
}

func TestEnsureConfigFileExists(t *testing.T) {
	t.Run("creates the .td directory and config file if it doesn't exist", func(t *testing.T) {
		tdCreated := false
		configCreated := false
		r := &RealReaderWriter{
			MkdirAllFunc: func(path string, perm os.FileMode) error {
				tdCreated = true
				return nil
			},
			StatFunc: func(name string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
				configCreated = true
				return nil
			},
		}
		err := r.EnsureConfigFileExists()

		assertNoError(t, err)

		assertTdCreated(t, tdCreated, true)
		assertTodosFileCreated(t, configCreated, true)
	})

	t.Run("does not create the todos file if it already exists", func(t *testing.T) {
		configCreated := false
		r := &RealReaderWriter{
			MkdirAllFunc: func(path string, perm os.FileMode) error {
				return nil
			},
			StatFunc: func(name string) (os.FileInfo, error) {
				return &MockFileInfo{}, os.ErrNotExist
			},
			WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
				configCreated = true
				return nil
			},
		}

		err := r.EnsureConfigFileExists()

		assertNoError(t, err)
		assertTodosFileCreated(t, configCreated, true)
	})

	t.Run("returns error when unable to create .td directory", func(t *testing.T) {
		r := &RealReaderWriter{
			MkdirAllFunc: func(path string, perm os.FileMode) error {
				return errors.New("unable to create .td directory")
			},
		}

		err := r.EnsureConfigFileExists()

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestReadConfigFile(t *testing.T) {
	t.Run("returns config from config.json", func(t *testing.T) {
		r := &RealReaderWriter{
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte(`{"hideCompleted":true}`), nil
			},
		}

		config, err := r.ReadConfigFile()

		assertNoError(t, err)
		if config.HideCompleted != true {
			t.Errorf("Expected hideCompleted to be true, got %v", config.HideCompleted)
		}
	})

	t.Run("returns error when read fails", func(t *testing.T) {
		var want error = errors.New("unable to read file")
		r := &RealReaderWriter{
			ReadFileFunc: func(filename string) ([]byte, error) {
				return nil, want
			},
		}

		_, got := r.ReadConfigFile()

		assertError(t, got, want)
	})

	t.Run("returns error when unmarshalling fails", func(t *testing.T) {
		r := &RealReaderWriter{
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte(`[ded]`), nil
			},
		}

		_, got := r.ReadConfigFile()

		if got == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func assertTdCreated(t *testing.T, gotTdCreated, wantTdCreated bool) {
	t.Helper()

	if gotTdCreated != wantTdCreated {
		t.Errorf("Expected tdCreated = %v, got %v", wantTdCreated, gotTdCreated)
	}
}

func assertTodosFileCreated(t *testing.T, gotTodosCreated, wantTodosCreated bool) {
	t.Helper()

	if gotTodosCreated != wantTodosCreated {
		t.Errorf("Expected todosCreated = %v, got %v", wantTodosCreated, gotTodosCreated)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got.Error() != want.Error() {
		t.Errorf("got \"%s\" want \"%s\"", got.Error(), want.Error())
	}
}
