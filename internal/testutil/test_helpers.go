package testutil

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var (
	// Warn is used when printing warnings
	Warn = colored(31)
	// Okay is used when printing text signifying a passed test
	Okay = colored(32)
	// Info is used for info messages
	Info = colored(35)
	// Notice is used for less important info messages
	Notice = colored(36)
	// String is used when printing strings
	String = colored(33)
)

func colored(code int) func(v ...interface{}) string {
	return func(v ...interface{}) string {
		return fmt.Sprint(append([]interface{}{shell(code)}, append(v, shell(0))...)...)
	}
}

func shell(i int) string {
	return "\x1B[" + strconv.Itoa(i) + "m"
}

func quoted(s string) string {
	var n string
	if len(strings.Split(s, "\n")) > 1 {
		n = "\n"
	}
	return fmt.Sprint(n, `"`, String(s), `"`)
}

// Compare data with the expected data and prints test output
func Compare(t *testing.T, expected string, out string) {
	if out != expected {
		t.Error(Warn(`FAIL!`), "\n\n"+
			`Got:      `+quoted(out)+"\n"+
			`Expected: `+quoted(expected)+"\n",
		)
	} else {
		t.Log(Okay(`PASS:`), quoted(strings.TrimSuffix(out, "\n")))
	}
}
