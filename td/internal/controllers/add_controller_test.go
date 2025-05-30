package controllers

import (
	"fmt"
	"os"
	"testing"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/testutils"
)

func TestAddController(t *testing.T) {
	t.Run("Run adds new todo", func(t *testing.T) {
		os.Setenv("DV_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		defer os.RemoveAll(dirname)

		controller := AddController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: false}, Args: []string{"do homework"}}, MetadataString: ""}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"do homework\" added to list.", got)
		testutils.FileContentEquals(t, fmt.Sprintf("- [ ] do homework %s\n", METADATA_STRING), file)
	})

	t.Run("Run correctly parses metadata string", func(t *testing.T) {
		os.Setenv("DV_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		defer os.RemoveAll(dirname)

		controller := AddController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: false}, Args: []string{"do homework"}}, MetadataString: "duration=10m,type=deep"}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"do homework\" added to list.", got)
		testutils.FileContentEquals(t, fmt.Sprintf("- [ ] do homework %s\n", "(duration=10m,type=deep)"), file)
	})
}
