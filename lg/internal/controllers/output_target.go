package controllers

import (
	"io"
	"os/exec"
	"strings"
)

type OutputTarget struct {
	Buffer strings.Builder
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
}

func (outputTarget *OutputTarget) Write(p []byte) (int, error) {
	if outputTarget.Stdin != nil {
		return outputTarget.Stdin.Write(p)
	} else {
		return outputTarget.Buffer.Write(p)
	}
}

func (outputTarget *OutputTarget) WriteString(p string) (int, error) {
	return outputTarget.Write([]byte(p))
}

func (outputTarget *OutputTarget) Close() error {
	var err error
	if outputTarget.Stdin != nil {
		err = outputTarget.Stdin.Close()
	}
	if outputTarget.Cmd != nil {
		err = outputTarget.Cmd.Wait()
	}
	return err
}

func (outputTarget *OutputTarget) String() string {
	return outputTarget.Buffer.String()
}

func NewCmdOutputTarget(writer io.WriteCloser, cmd *exec.Cmd) OutputTarget {
	return OutputTarget{Stdin: writer, Cmd: cmd, Buffer: strings.Builder{}}
}

func NewSimpleOutputTarget() OutputTarget {
	return OutputTarget{Stdin: nil, Cmd: nil, Buffer: strings.Builder{}}
}
