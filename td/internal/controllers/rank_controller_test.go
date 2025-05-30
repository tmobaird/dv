package controllers

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/testutils"
)

func TestRankController(t *testing.T) {
	t.Run("Run ranks todo correctly when hideCompleted false", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two\n- [ ] Three\n"))
		defer os.RemoveAll(dirname)

		controller := RankController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: false}, Args: []string{"1", "2"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"One\" ranked to index 2.", got)

		_, err = file.Seek(0, 0)
		testutils.AssertNoError(t, err)
		data, err := io.ReadAll(file)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, fmt.Sprintf("- [ ] Two %s\n- [x] One %s\n- [ ] Three %s\n", METADATA_STRING, METADATA_STRING, METADATA_STRING), string(data))
	})

	t.Run("Run ranks todo correctly hideCompleted true", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two\n- [ ] Three\n"))
		defer os.RemoveAll(dirname)

		controller := RankController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: true}, Args: []string{"1", "2"}}}
		got, err := controller.Run()

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, "\"Two\" ranked to index 2.", got)

		_, err = file.Seek(0, 0)
		testutils.AssertNoError(t, err)
		data, err := io.ReadAll(file)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, fmt.Sprintf("- [x] One %s\n- [ ] Three %s\n- [ ] Two %s\n", METADATA_STRING, METADATA_STRING, METADATA_STRING), string(data))
	})

	t.Run("Run returns an error when target index is out of range", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two\n- [ ] Three\n"))
		defer os.RemoveAll(dirname)

		controller := RankController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: true}, Args: []string{"1000", "2"}}}
		_, err = controller.Run()

		testutils.AssertEqual(t, "todo @ index 1000 does not exist", err.Error())
	})

	t.Run("Run returns an error when end index is out of range", func(t *testing.T) {
		os.Setenv("TD_BASE_PATH", "tmp")
		dirname := "tmp/lists"
		err := os.MkdirAll(dirname, 0755)
		testutils.AssertNoError(t, err)
		file, err := os.Create("tmp/lists/main.md")
		testutils.AssertNoError(t, err)
		file.Write([]byte("- [x] One\n- [ ] Two\n- [ ] Three\n"))
		defer os.RemoveAll(dirname)

		controller := RankController{Base: Controller{Config: core.Config{Context: "main", HideCompleted: true}, Args: []string{"1", "2000"}}}
		_, err = controller.Run()

		testutils.AssertEqual(t, "cant move todo to @2000, out of bounds", err.Error())
	})
}
