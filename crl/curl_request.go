package crl

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"
)

type CurlRequest struct {
	Url     *url.URL
	Method  RequestMethod
	Headers map[string]string
	Body    string
}

func quoteArg(arg string) string {
	if strings.ContainsAny(arg, " \t\n\"'`$&|<>();") {
		return "'" + strings.ReplaceAll(arg, "'", `'"'"'`) + "'"
	}
	return arg
}

// BuildCommandString reconstructs the shell-equivalent command
func buildCommandString(name string, args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = quoteArg(arg)
	}
	return name + " " + strings.Join(quoted, " ")
}

func (request CurlRequest) Execute(logger io.Writer) ([]byte, error) {
	args := []string{}

	args = append(args, "--request")
	args = append(args, request.Method.toString())
	for key, value := range request.Headers {
		args = append(args, "-H")
		args = append(args, fmt.Sprintf("%s: %s", key, value))
	}
	fmt.Println(request.Url, request.Url.String())
	args = append(args, request.Url.String())

	logger.Write([]byte(fmt.Sprintf("Making curl request: %s\n\n", buildCommandString("curl", args))))
	command := exec.Command("curl", args...)
	output := bytes.Buffer{}
	command.Stdout = &output

	err := command.Run()
	return output.Bytes(), err
}
