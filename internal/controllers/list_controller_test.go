package controllers

import (
	"os"
	"td/internal"
	"td/internal/testutils"
	"testing"
)

func TestListController(t *testing.T) {
	t.Run("Run returns list of todos when exist", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := ListController{Base: Controller{Config: internal.DefaultConfig(), Args: []string{}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "1. [x] One\n2. [ ] Two\n", got)
	})

	t.Run("Run returns list that doesnt include misformatted lines", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- messed up"))
		defer os.RemoveAll(dirname)

		controller := ListController{Base: Controller{Config: internal.DefaultConfig(), Args: []string{}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "No todos in list.", got)
	})

	t.Run("Run returns list with completed only when Config.HideCompleted = true", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two"))
		defer os.RemoveAll(dirname)

		controller := ListController{Base: Controller{Config: internal.Config{Context: "main", HideCompleted: true}, Args: []string{}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "1. [ ] Two\n", got)
	})
}
