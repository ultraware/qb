package qb

import (
	"strconv"
)

// Values used when building queries
var (
	COMMA   = `, `
	NEWLINE = "\n"
	INDENT  = "\t"
	VALUE   = `?`
)

///// Field /////

// MakeField returns the value as a Field, no operation performed when the value is already a field
func MakeField(i interface{}) Field {
	if f, ok := i.(Field); ok {
		return f
	}
	return Value(i)
}

// ConcatQuery combines strings and QueryStringers into string
func ConcatQuery(c *Context, values ...interface{}) string {
	s := ``
	for _, val := range values {
		switch v := val.(type) {
		case (Field):
			s += v.QueryString(c)
		case (string):
			s += v
		}
	}
	return s
}

// JoinQuery joins fields or values into a string separated by sep
func JoinQuery(c *Context, sep string, values []interface{}) string {
	s := make([]interface{}, len(values)*2-1)
	for k, v := range values {
		if k > 0 {
			s[k*2-1] = sep
		}
		s[k*2] = MakeField(v)
	}

	return ConcatQuery(c, s...)
}

///// Alias /////

type noAlias struct{}

func (n *noAlias) Get(_ Source) string {
	return ``
}

// NoAlias returns no alias
func NoAlias() Alias {
	return &noAlias{}
}

type aliasGenerator struct {
	cache    map[Source]string
	counters map[string]int
}

// AliasGenerator returns an incrementing alias for each new Source
func AliasGenerator() Alias {
	return &aliasGenerator{make(map[Source]string), make(map[string]int)}
}

func (g *aliasGenerator) Get(src Source) string {
	if src == nil {
		return ``
	}

	if v, ok := g.cache[src]; ok {
		return v
	}

	new := g.new(src)
	g.cache[src] = new
	return new
}

func (g *aliasGenerator) new(src Source) string {
	a := src.aliasString()

	g.counters[a]++

	if g.counters[a] == 1 {
		return a
	}

	return a + strconv.Itoa(g.counters[a])
}
