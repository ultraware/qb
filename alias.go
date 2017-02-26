package qb

import "strconv"

type aliasGenerator struct {
	counter int
	list    map[string]string
}

func newGenerator() aliasGenerator {
	return aliasGenerator{0, make(map[string]string)}
}

func (g *aliasGenerator) Get(src Source) string {
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
