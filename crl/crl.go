package crl

import (
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

type ArgMetadata struct {
	Url         string
	Method      string
	ContentType string
}

func createArgMetadata(args ...string) ArgMetadata {
	url := "http://replace"
	method := "GET"

	rawArgs := strings.Join(args, " ")

	// parse url
	urlRegex := regexp.MustCompile(`--url=(\S+)`)
	urlMatches := urlRegex.FindStringSubmatch(rawArgs)
	if len(urlMatches) > 1 {
		url = urlMatches[1]
	}

	// parse method
	methodRegex := regexp.MustCompile(`--method=(GET|POST|PUT|PATCH|DELETE)`)
	methodMatches := methodRegex.FindStringSubmatch(rawArgs)
	if len(methodMatches) > 1 {
		method = methodMatches[1]
	}

	contentType := "application/json"
	if strings.Contains(rawArgs, "--form") {
		contentType = "application/x-www-form-urlencoded"
	} else if strings.Contains(rawArgs, "plain") {
		contentType = "text/plain"
	} else if strings.Contains(rawArgs, "xml") {
		contentType = "application/xml"
	}
	return ArgMetadata{
		Url:         url,
		Method:      method,
		ContentType: contentType,
	}
}

func copyTemplateContentsToTemp(file *os.File, argMetadata ArgMetadata) error {
	_, filename, _, _ := runtime.Caller(0)
	tmplPath := filepath.Join(filepath.Dir(filename), "template.md")
	tmpl, err := template.ParseFiles(tmplPath)

	if err != nil {
		return err
	}
	err = tmpl.Execute(file, argMetadata)
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
	groups = groups[1:] // we dont care about the first item, cause it should be blank
	if len(groups) < 2 {
		return request, errors.New("must include more parts")
	}

	metadata, err := parseFrontMatter(groups[0])
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
	request.Body = parseBody(groups[1])
	return request, nil
}

var crlCmd = &cobra.Command{
	Use:                "crl",
	Short:              "a curl enhancement",
	Long:               "TODO",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// create tmp file
		file, err := core.CreateTempFile("*.md")
		if err != nil {
			writeErr(cmd, err)
		}
		defer os.Remove(file.Name())

		// get metadata from args
		metadata := createArgMetadata(os.Args...)

		// copy contents from template
		err = copyTemplateContentsToTemp(file, metadata)
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

func parseFrontMatter(frontMatter string) (map[string]interface{}, error) {
	frontMatter = strings.Trim(frontMatter, "\n")
	var metadata map[string]interface{}
	err := yaml.Unmarshal([]byte(frontMatter), &metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
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

func parseBody(rawBody string) string {
	rawBody = strings.Trim(rawBody, "\n")
	if strings.HasPrefix(rawBody, "```") {
		lines := strings.Split(rawBody, "\n")
		lines = lines[1:] // throw away line one
		if strings.HasPrefix(lines[len(lines)-1], "```") {
			lines = lines[0 : len(lines)-1]
		}
		return strings.Join(lines, "\n")
	} else {
		return rawBody
	}
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
