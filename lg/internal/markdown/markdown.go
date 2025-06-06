package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tmobaird/dv/colors"
)

func MdToOutput(markdown string) string {
	lines := strings.Split(markdown, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = substituteHeaders(lines[i])
		lines[i] = substituteInlineStyles(lines[i])
	}
	return strings.Join(lines, "\n")
}

func substituteHeaders(line string) string {
	headerRegex := regexp.MustCompile(`^(#{1,6}) (.+)`)
	groups := headerRegex.FindStringSubmatch(line)
	newLine := ""
	if len(groups) == 3 {
		switch len(groups[1]) {
		case 1:
			newLine = colors.AddTextStyle(colors.AddTextStyle(groups[2], colors.CODE_UNDERLINE), colors.CODE_BOLD)
		case 2:
			newLine = colors.AddTextStyle(groups[2], colors.CODE_UNDERLINE)
		default:
			newLine = colors.AddTextStyle(groups[2], colors.CODE_BOLD)
		}
		return newLine
	} else {
		return line
	}
}

func substituteInlineStyles(line string) string {
	boldRegex := regexp.MustCompile(`(?:[\*]{2}|[_]{2})(.+)(?:[\*]{2}|[_]{2})`)
	groups := boldRegex.FindStringSubmatch(line)
	if len(groups) > 1 {
		for _, match := range groups[1:] {
			start := strings.Index(line, fmt.Sprintf("**%s**", match))
			if start == -1 {
				start = strings.Index(line, fmt.Sprintf("__%s__", match))
			}
			length := len(match) + 4
			chars := strings.Split(line, "")
			line = strings.Join(chars[0:start], "") + colors.AddTextStyle(match, colors.CODE_BOLD) + strings.Join(chars[start+length:], "")
		}
	}

	italicsRegex := regexp.MustCompile(`(?:[\*]{1}|[_]{1})(.+)(?:[\*]{1}|[_]{1})`)
	groups = italicsRegex.FindStringSubmatch(line)
	if len(groups) > 1 {
		for _, match := range groups[1:] {
			start := strings.Index(line, fmt.Sprintf("*%s*", match))
			if start == -1 {
				start = strings.Index(line, fmt.Sprintf("_%s_", match))
			}
			length := len(match) + 2
			chars := strings.Split(line, "")
			line = strings.Join(chars[0:start], "") + colors.AddTextStyle(match, colors.CODE_ITALICS) + strings.Join(chars[start+length:], "")
		}
	}

	return line
}
