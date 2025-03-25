package qb

import "strconv"

type noAlias struct{}

func (n *noAlias) Get(_ Source) string {
	return ``
}

// NoAlias returns no alias.
// This function is not intended to be called directly
func NoAlias() Alias {
	return &noAlias{}
}

type aliasGenerator struct {
	cache    map[Source]string
	counters map[string]int
}

// AliasGenerator returns an incrementing alias for each new Source.
// This function is not intended to be called directly
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

	nv := g.new(src)
	g.cache[src] = nv
	return nv
}

func (g *aliasGenerator) new(src Source) string {
	a := src.aliasString()

	g.counters[a]++

	if g.counters[a] == 1 {
		return a
	}

	return a + strconv.Itoa(g.counters[a])
}
