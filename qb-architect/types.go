package main

import (
	"sort"

	"github.com/Ultraware/qb/v3/qb-architect/internal/db"
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
