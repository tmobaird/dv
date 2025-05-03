package testutils

import (
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
