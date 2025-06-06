package markdown

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tmobaird/dv/colors"
	"github.com/tmobaird/dv/testutils"
)

func TestMarkdown(t *testing.T) {
	t.Run("1 Level Headers are appropriately substituted", func(t *testing.T) {
		expected := fmt.Sprintf("%s%sHello%s%s", colors.EscapeCode(colors.CODE_BOLD), colors.EscapeCode(colors.CODE_UNDERLINE), colors.EscapeCode(colors.CODE_RESET_UNDERLINE), colors.EscapeCode(colors.CODE_RESET_BOLD))
		got := MdToOutput("# Hello")

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("2 Level Headers are appropriately substituted", func(t *testing.T) {
		expected := fmt.Sprintf("%sHello%s", colors.EscapeCode(colors.CODE_UNDERLINE), colors.EscapeCode(colors.CODE_RESET_UNDERLINE))
		got := MdToOutput("## Hello")

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("All other level headers are appropriately substituted", func(t *testing.T) {
		for i := 3; i < 7; i++ {
			expected := fmt.Sprintf("%sHello%s", colors.EscapeCode(colors.CODE_BOLD), colors.EscapeCode(colors.CODE_RESET_BOLD))
			got := MdToOutput(fmt.Sprintf("%s Hello", strings.Repeat("#", i)))

			testutils.AssertEqual(t, expected, got)
		}
	})

	t.Run("'# ' in the middle of a line does nothing", func(t *testing.T) {
		testutils.AssertEqual(t, "hello # world", MdToOutput("hello # world"))
	})

	t.Run("**Bold** characters properly get substituted", func(t *testing.T) {
		expected := fmt.Sprintf("%sHello%s", colors.EscapeCode(colors.CODE_BOLD), colors.EscapeCode(colors.CODE_RESET_BOLD))
		got := MdToOutput("**Hello**")

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("__Bold__ characters properly get substituted", func(t *testing.T) {
		expected := fmt.Sprintf("%sHello%s", colors.EscapeCode(colors.CODE_BOLD), colors.EscapeCode(colors.CODE_RESET_BOLD))
		got := MdToOutput("__Hello__")

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("*Italics* characters properly get substituted", func(t *testing.T) {
		expected := fmt.Sprintf("%sHello%s", colors.EscapeCode(colors.CODE_ITALICS), colors.EscapeCode(colors.CODE_RESET_ITALICS))
		got := MdToOutput("*Hello*")

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("_Italics_ characters properly get substituted", func(t *testing.T) {
		expected := fmt.Sprintf("%sHello%s", colors.EscapeCode(colors.CODE_ITALICS), colors.EscapeCode(colors.CODE_RESET_ITALICS))
		got := MdToOutput("_Hello_")

		testutils.AssertEqual(t, expected, got)
	})

	t.Run("Italics and Bold can occur on same line", func(t *testing.T) {
		expected := fmt.Sprintf("My %sname%s is %sthomas%s", colors.EscapeCode(colors.CODE_BOLD), colors.EscapeCode(colors.CODE_RESET_BOLD), colors.EscapeCode(colors.CODE_ITALICS), colors.EscapeCode(colors.CODE_RESET_ITALICS))
		got := MdToOutput("My **name** is _thomas_")

		testutils.AssertEqual(t, expected, got)
	})
}
