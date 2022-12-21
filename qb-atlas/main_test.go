package main

import "testing"

func TestMain(m *testing.M) {
	m.Run()
}

func eq[T string | int | bool](t *testing.T, expected, actual T) {
	if actual != expected {
		t.Errorf("expected: '%v', but got: '%v'", expected, actual)
	}
}

func noErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func hasErr(t *testing.T, err error) {
	if err == nil {
		t.Error("missing expected error")
	}
}
