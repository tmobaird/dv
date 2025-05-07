package internal

import (
	"bytes"
	"os"
	"strings"
	"td/internal/testutils"
	"testing"
	"testing/fstest"
)

func TestConfig(t *testing.T) {
	t.Run("FileExists returns true when file exists", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"config.yaml": {Data: []byte("context: main")},
		}

		expected := true
		got := FileExists(filesystem)

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("FileExists returns false when file does not exist", func(t *testing.T) {
		filesystem := fstest.MapFS{}

		expected := false
		got := FileExists(filesystem)

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("Read returns existing config from file when exists", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"config.yaml": {Data: []byte("context: my-sandbox")},
		}

		expected := "my-sandbox"
		config, e := Read(filesystem)

		testutils.AssertNoError(t, e)
		testutils.AssertEqual(t, expected, config.Context)
	})

	t.Run("Read returns error when fails to unmarshal", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"config.yaml": {Data: []byte("fdasfjdsklajfldks")},
		}

		_, e := Read(filesystem)

		testutils.AssertError(t, e)
	})

	t.Run("Read returns default config when file does not exist", func(t *testing.T) {
		filesystem := fstest.MapFS{}

		expected := "main"
		config, e := Read(filesystem)

		testutils.AssertNoError(t, e)
		testutils.AssertEqual(t, expected, config.Context)
	})

	t.Run("Save writes the file", func(t *testing.T) {
		output := &bytes.Buffer{}
		config := Config{Context: "dummy"}

		err := Save(output, config)
		testutils.AssertNoError(t, err)

		written := output.String()
		if !strings.Contains(written, "context: dummy") {
			t.Errorf("Expected %s, to include \"context: dummy\"", written)
		}
	})

	t.Run("PersistConfig writes config to file", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		file, err := os.Create("tmp/config.yaml")
		if err != nil {
			t.Errorf("Failed to create tmp file %s", err.Error())
		}
		defer os.Remove(file.Name())
		config := Config{Context: "my-context"}
		Save(file, config)

		config.Context = "updated-context"
		PersistConfig(config)

		config, err = Read(os.DirFS(BasePath()))
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, "updated-context", config.Context)
	})

	t.Run("PersistConfig does not return an error when having to create the file", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		defer os.Remove("tmp/config.yaml")
		config := Config{Context: "my-context"}
		err := PersistConfig(config)
		testutils.AssertNoError(t, err)
	})
}
