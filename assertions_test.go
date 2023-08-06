package main

import (
	"testing"
)

func assertEquals(t testing.TB, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got \"%v\" want \"%v\"", got, want)
	}
}

func assertCorrectError(t testing.TB, got error, want string) {
	t.Helper()
	if got.Error() != want {
		t.Errorf("got \"%s\" want \"%s\"", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
