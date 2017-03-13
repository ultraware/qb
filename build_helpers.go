package qb

import "strconv"

// AliasGenerator makes aliasses for tables and keeps track of the previously given aliasses
type AliasGenerator struct {
	counter int
	list    map[string]string
}

func newGenerator() AliasGenerator {
	return AliasGenerator{0, make(map[string]string)}
}

// Get returns the alias for the given source
func (g *AliasGenerator) Get(src Source) string {
	if src == nil {
		return ``
	}

	if v, ok := g.list[src.QueryString()]; ok {
		return v
	}

	g.counter++
	g.list[src.QueryString()] = `a` + strconv.Itoa(g.counter)
	return g.list[src.QueryString()]
}

// ValueList is a list of static values used in a query
type ValueList []interface{}

// Append adds the given values to the list
func (list *ValueList) Append(v ...interface{}) {
	*list = append(*list, v...)
}
