package filter

import (
	"fmt"
	"os"
	"regexp"
)

// Filters is a list of regular expressions
type Filters []*regexp.Regexp

// Set implements flag.Value
func (f *Filters) Set(value string) error {
	re, err := regexp.Compile(value)
	if err != nil {
		println(`Invalid regular expression: ` + value + "\n" + err.Error())
		os.Exit(2)
	}

	*f = append(*f, re)

	return nil
}

// String implements flag.Value
func (f Filters) String() string {
	return fmt.Sprint([]*regexp.Regexp(f))
}
