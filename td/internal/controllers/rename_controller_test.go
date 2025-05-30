package controllers

import (
	"os"
	"testing"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/testutils"
)

func TestRenameController(t *testing.T) {
	t.Run("Run renames correct todo when hideCompleted false", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := RenameController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: false}, Args: []string{"1", "Updated Name"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"One\" updated to \"Updated Name\".", got)
	})

	t.Run("Run renames correct todo when hideCompleted true", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := RenameController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: true}, Args: []string{"1", "Updated Name"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"Two\" updated to \"Updated Name\".", got)
	})

	t.Run("Run renames an error when index is out of range", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte(""))
		defer os.RemoveAll(dirname)

		controller := RenameController{Base: Controller{Config: core.DefaultConfig(), Args: []string{"1000", "Updated Name"}}}
		_, err = controller.Run()

		testutils.AssertError(t, err)
		testutils.AssertEqual(t, "todo @ index 1000 does not exist", err.Error())
	})
}
