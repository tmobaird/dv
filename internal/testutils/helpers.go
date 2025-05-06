package testutils

import (
	"io"
	"os"
	"testing"
)

func AssertEqual[T comparable](t *testing.T, expected, got T, msg ...interface{}) {
	t.Helper()
	if expected != got {
		t.Errorf("\nexpected '%v'\ngot '%v'\n%v", expected, got, msg)
	}
}

func AssertNotEqual[T comparable](t *testing.T, expected, got T, msg ...interface{}) {
	t.Helper()
	if expected == got {
		t.Errorf("expected '%v'\nto not equal '%v'\n%v", got, expected, msg)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("did not expect an error, got %s", err.Error())
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("expected an error")
	}
}

func FileContentEquals(t *testing.T, expected string, file *os.File) {
	t.Helper()

	_, err := file.Seek(0, 0)
	AssertNoError(t, err)
	data, err := io.ReadAll(file)
	AssertNoError(t, err)
	AssertEqual(t, expected, string(data))
}
