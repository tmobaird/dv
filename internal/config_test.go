package internal

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"
)

func AssertEqual[T comparable](t *testing.T, expected, got T, msg ...interface{}) {
	t.Helper()
	if expected != got {
		t.Errorf("got %v, want %v, %v", expected, got, msg)
	}
}

func AssertNotEqual[T comparable](t *testing.T, expected, got T, msg ...interface{}) {
	t.Helper()
	if expected == got {
		t.Errorf("did not want %v, got %v, %v", got, expected, msg)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("did not expect an error, got %s", err.Error())
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("expected an error")
	}
}

func TestConfig(t *testing.T) {
	t.Run("FileExists returns true when file exists", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"config.yaml": {Data: []byte("context: main")},
		}

		expected := true
		got := FileExists(filesystem)

		AssertEqual(t, expected, got)
	})

	t.Run("FileExists returns false when file does not exist", func(t *testing.T) {
		filesystem := fstest.MapFS{}

		expected := false
		got := FileExists(filesystem)

		AssertEqual(t, expected, got)
	})

	t.Run("Read returns existing config from file when exists", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"config.yaml": {Data: []byte("context: my-sandbox")},
		}

		expected := "my-sandbox"
		config, e := Read(filesystem)

		AssertNoError(t, e)
		AssertEqual(t, expected, config.Context)
	})

	t.Run("Read returns error when fails to unmarshal", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"config.yaml": {Data: []byte("fdasfjdsklajfldks")},
		}

		_, e := Read(filesystem)

		AssertError(t, e)
	})

	t.Run("Read returns default config when file does not exist", func(t *testing.T) {
		filesystem := fstest.MapFS{}

		expected := "main"
		config, e := Read(filesystem)

		AssertNoError(t, e)
		AssertEqual(t, expected, config.Context)
	})

	t.Run("Save writes the file", func(t *testing.T) {
		output := &bytes.Buffer{}
		config := Config{Context: "dummy"}

		err := Save(output, config)
		AssertNoError(t, err)

		written := output.String()
		if !strings.Contains(written, "context: dummy") {
			t.Errorf("Expected %s, to include \"context: dummy\"", written)
		}
	})
}
