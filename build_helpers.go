package qb

import "strconv"

// Constants used when building queries
const (
	COMMA = `, `
	VALUE = `?`
)

// NoAlias returns no alias
type NoAlias struct{}

// Get implements the Alias interface
func (n *NoAlias) Get(_ Source) string {
	return ``
}

// AliasGenerator makes aliasses for tables and keeps track of the previously given aliasses
type AliasGenerator struct {
	counter int
	list    map[Source]string
}

func newGenerator() AliasGenerator {
	return AliasGenerator{0, make(map[Source]string)}
}

// Get returns the alias for the given source
func (g *AliasGenerator) Get(src Source) string {
	if src == nil {
		return ``
	}

	if v, ok := g.list[src]; ok {
		return v
	}

	g.counter++
	g.list[src] = src.AliasString() + strconv.Itoa(g.counter)
	return g.list[src]
}

// ValueList is a list of static values used in a query
type ValueList []interface{}

// Append adds the given values to the list
func (list *ValueList) Append(v ...interface{}) {
	*list = append(*list, v...)
}
