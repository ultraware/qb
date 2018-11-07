package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kortschak/ct"
)

var (
	// Warn is used when printing warnings
	Warn = (ct.Fg(ct.White) | ct.Bg(ct.Red) | ct.Bold).Paint
	// Okay is used when printing text signifying a passed test
	Okay = (ct.Fg(ct.Green) | ct.Bold).Paint
	// Info is used for info messages
	Info = (ct.Fg(ct.Magenta)).Paint
	// Notice is used for less important info messages
	Notice = (ct.Fg(ct.Magenta) | ct.Faint).Paint

	// String is used when printing strings
	String = (ct.Fg(ct.Yellow).Paint)
)

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
