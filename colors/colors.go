package colors

import "fmt"

const ESCAPE = "\x1b"

const CODE_RESET = 0

// Text Styles
const CODE_BOLD = 1
const CODE_UNDERLINE = 4
const CODE_RESET_BOLD = 22
const CODE_RESET_UNDERLINE = 24

var resetMap map[int]int = map[int]int{
	CODE_BOLD:      CODE_RESET_BOLD,
	CODE_UNDERLINE: CODE_RESET_UNDERLINE,
}

// Foreground Colors
const FG_BLACK = 30
const FG_RED = 31
const FG_GREEN = 32
const FG_YELLOW = 33
const FG_BLUE = 34
const FG_MAGENTA = 35
const FG_CYAN = 36
const FG_WHITE = 37

func AddColor(str string, color int) string {
	return fmt.Sprintf("%s%s%s", EscapeCode(color), str, EscapeCode(CODE_RESET))
}

func AddTextStyle(str string, style int) string {
	resetCode := resetMap[style]
	return fmt.Sprintf("%s%s%s", EscapeCode(style), str, EscapeCode(resetCode))
}

func EscapeCode(code int) string {
	return fmt.Sprintf("%s[%dm", ESCAPE, code)
}
