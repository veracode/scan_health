package utils

import "testing"

func nope(t *testing.T, message string) {
	t.Errorf("%s", message)
}

func expectString(t *testing.T, result, expected string) {
	if result != expected {
		t.Errorf("got:\n%q, wanted:\n%q", result, expected)
	}
}
