package controllers

import (
	"os"
	"td/internal"
	"td/internal/testutils"
	"testing"
)

func TestMarkController(t *testing.T) {
	t.Run("Run marks todo as done when hideCompleted is false", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: false}, Args: []string{"2", "done"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"Two\" marked done.", got)
		testutils.FileContentEquals(t, "- [x] One\n- [x] Two\n", file)
	})

	t.Run("Run marks todo as done when hideCompleted is true", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: true}, Args: []string{"1", "done"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"Two\" marked done.", got)
		testutils.FileContentEquals(t, "- [x] One\n- [x] Two\n", file)
	})

	t.Run("Run marks todo as done when given d", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [ ] One"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: false}, Args: []string{"1", "d"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"One\" marked done.", got)
		testutils.FileContentEquals(t, "- [x] One\n", file)
	})

	t.Run("Run marks todo as not done when given not-done", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: false}, Args: []string{"1", "not-done"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"One\" marked not done.", got)

		testutils.FileContentEquals(t, "- [ ] One\n", file)
	})

	t.Run("Run marks todo as not done when given not-done", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: false}, Args: []string{"1", "not"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"One\" marked not done.", got)
	})

	t.Run("Run with bad index returns error", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: false}, Args: []string{"1000", "not"}}}
		_, err = controller.Run()

		testutils.AssertError(t, err)
		testutils.AssertEqual(t, "todo @ index 1000 does not exist", err.Error())
	})

	t.Run("Run with bad status returns error", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One"))
		defer os.RemoveAll(dirname)

		controller := MarkController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: false}, Args: []string{"1", "bad"}}}
		_, err = controller.Run()

		testutils.AssertError(t, err)
		testutils.AssertEqual(t, "status must be one of: [done, d, not-done, not], got bad", err.Error())
	})
}
