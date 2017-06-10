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
func ConcatQuery(d Driver, ag Alias, vl *ValueList, values ...interface{}) string {
	s := ``
	for _, val := range values {
		switch v := val.(type) {
		case (Field):
			s += v.QueryString(d, ag, vl)
		case (string):
			s += v
		}
	}
	return s
}

///// Alias /////

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

func newGenerator() *AliasGenerator {
	return &AliasGenerator{0, make(map[Source]string)}
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

///// Value list /////

// ValueList is a list of static values used in a query
type ValueList []interface{}

// Append adds the given values to the list
func (list *ValueList) Append(v ...interface{}) {
	*list = append(*list, v...)
}

///// Primary /////

// GetPrimaryFields return all fields in the given list that are a primary key
func GetPrimaryFields(f []DataField) []DataField {
	list := []DataField{}
	for _, v := range f {
		if !isPrimary(v) {
			continue
		}
		list = append(list, v)
	}
	return list
}

func isPrimary(v DataField) bool {
	f, ok := v.Field.(*TableField)
	return ok && f.Primary
}
