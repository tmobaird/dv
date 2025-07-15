package crl

import (
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/tmobaird/dv/core"
	"gopkg.in/yaml.v3"
)

func (method RequestMethod) toString() string {
	switch method {
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Patch:
		return "PATCH"
	case Delete:
		return "DELETE"
	default:
		return "GET"
	}
}

type RequestMethod int

const (
	Get RequestMethod = iota
	Post
	Put
	Patch
	Delete
)

func copyTemplateContentsToTemp(file *os.File) error {
	_, filename, _, _ := runtime.Caller(0)
	tmplPath := filepath.Join(filepath.Dir(filename), "template")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, nil)
	if err != nil {
		return err
	}
	return nil
}

func openFileInEditor(file *os.File) error {
	editor, err := core.Editor()
	if err != nil {
		return err
	}
	command := exec.Command(editor, file.Name())
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	err = command.Run()
	return err
}

func parseTempFile(file *os.File) (CurlRequest, error) {
	request := CurlRequest{}
	// execute curl request on exit
	contents, err := os.ReadFile(file.Name())
	if err != nil {
		return request, err
	}

	// parse the file contents
	fileContents := string(contents)
	groups := strings.Split(fileContents, "---")
	if len(groups) < 3 {
		return request, errors.New("must include more parts")
	}

	// parse front matter
	rawMetadata := strings.Trim(groups[1], "\n")
	var metadata map[string]interface{}
	err = yaml.Unmarshal([]byte(rawMetadata), &metadata)
	if err != nil {
		return request, err
	}

	// create request
	url, err := parseUrl(metadata)
	if err != nil {
		return request, err
	}
	request.Url = url
	request.Method = parseMethod(metadata)
	request.Headers = parseHeaders(metadata)
	request.Body = groups[2]
	return request, nil
}

var crlCmd = &cobra.Command{
	Use:   "crl",
	Short: "a curl enhancement",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		// create tmp file
		file, err := core.CreateTempFile("")
		if err != nil {
			writeErr(cmd, err)
		}
		defer os.Remove(file.Name())

		// copy contents from template
		err = copyTemplateContentsToTemp(file)
		if err != nil {
			writeErr(cmd, err)
			return
		}

		// open tmp file in editor
		err = openFileInEditor(file)
		if err != nil {
			writeErr(cmd, err)
			return
		}

		// execute curl request on exit
		request, err := parseTempFile(file)
		if err != nil {
			writeErr(cmd, err)
			return
		}

		output, err := request.Execute(cmd.OutOrStdout())

		if err != nil {
			writeErr(cmd, err)
		} else {
			writeOut(cmd, output)
		}
	},
}

func parseUrl(metadata map[string]interface{}) (*url.URL, error) {
	raw, ok := metadata["url"].(string)
	if !ok {
		return &url.URL{}, errors.New("url must be a valid string")
	}

	return url.Parse(raw)
}

func parseMethod(metadata map[string]interface{}) RequestMethod {
	method, ok := metadata["method"].(string)
	if !ok {
		return Get
	}

	if slices.Contains([]string{"GET", "Get", "get"}, method) {
		return Get
	} else if slices.Contains([]string{"POST", "Post", "post"}, method) {
		return Post
	} else if slices.Contains([]string{"PUT", "Put", "put"}, method) {
		return Put
	} else if slices.Contains([]string{"PATCH", "Patch", "patch"}, method) {
		return Patch
	} else if slices.Contains([]string{"DELETE", "Delete", "delete"}, method) {
		return Delete
	}

	return Get
}

func parseHeaders(metadata map[string]interface{}) map[string]string {
	if headersRaw, ok := metadata["headers"]; ok {
		if headersList, ok := headersRaw.([]interface{}); ok && len(headersList) > 0 {
			// Now extract the first element and assert it is a map[string]interface{}
			if headerMap, ok := headersList[0].(map[string]interface{}); ok {
				result := make(map[string]string)
				for k, v := range headerMap {
					if strVal, ok := v.(string); ok {
						result[k] = strVal
					}
				}
				return result
			}
		}
	}
	return map[string]string{}
}

func writeErr(cmd *cobra.Command, err error) {
	cmd.OutOrStderr().Write([]byte(err.Error()))
}

func writeOut(cmd *cobra.Command, bytes []byte) {
	cmd.OutOrStdout().Write(bytes)
}

func CrlCommand() *cobra.Command {
	return crlCmd
}
