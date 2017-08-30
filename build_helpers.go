package qb

import "strconv"

// Values used when building queries
var (
	COMMA   = `, `
	NEWLINE = "\n"
	INDENT  = "\t"
	VALUE   = `?`
)

///// Field /////

// MakeField ...
func MakeField(i interface{}) Field {
	if f, ok := i.(Field); ok {
		return f
	}
	return Value(i)
}

// ConcatQuery ...
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
	counter int
	list    map[Source]string
}

// AliasGenerator returns an incrementing alias for each new Source
func AliasGenerator() Alias {
	return &aliasGenerator{0, make(map[Source]string)}
}

func (g *aliasGenerator) Get(src Source) string {
	if src == nil {
		return ``
	}

	if v, ok := g.list[src]; ok {
		return v
	}

	g.counter++
	g.list[src] = src.aliasString() + strconv.Itoa(g.counter)
	return g.list[src]
}
