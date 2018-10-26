package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"git.ultraware.nl/NiseVoid/qb/qb-architect/db"
)

type fields []db.Field

func (f fields) Len() int {
	return len(f)
}

func (f fields) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{f[i].Name, f[j].Name})
}

func (f fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type filters []*regexp.Regexp

func (f *filters) Set(value string) error {
	re, err := regexp.Compile(value)
	if err != nil {
		println(`Invalid regular expression: ` + value + "\n" + err.Error())
		os.Exit(2)
	}

	*f = append(*f, re)

	return nil
}

func (f filters) String() string {
	return fmt.Sprint([]*regexp.Regexp(f))
}
