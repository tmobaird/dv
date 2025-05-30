package controllers

import (
	"os"
	"testing"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/testutils"
)

func TestRemoveController(t *testing.T) {
	t.Run("Run deletes correct todo when hideCompleted false", func(t *testing.T) {
		os.Setenv("DV_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := RemoveController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: false}, Args: []string{"1"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"One\" removed from list.", got)
	})

	t.Run("Run deletes correct todo when hideCompleted true", func(t *testing.T) {
		os.Setenv("DV_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := RemoveController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: true}, Args: []string{"1"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"Two\" removed from list.", got)
	})

	t.Run("Run returns an error when index is out of range", func(t *testing.T) {
		os.Setenv("DV_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte(""))
		defer os.RemoveAll(dirname)

		controller := RemoveController{Base: Controller{Config: core.DefaultConfig(), Args: []string{"1000"}}}
		_, err = controller.Run()

		testutils.AssertError(t, err)
		testutils.AssertEqual(t, "todo @ index 1000 does not exist", err.Error())
	})
}
