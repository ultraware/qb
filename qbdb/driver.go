package qbdb

import (
	"strconv"

	"git.ultraware.nl/NiseVoid/qb"
)

// Driver is a default driver used for tests
type Driver struct{}

// ValueString returns the placeholder for prepare values
func (d Driver) ValueString(c int) string {
	return `@@`
}

// BoolString returns the notation for boolean values
func (d Driver) BoolString(v bool) string {
	if v {
		return `t`
	}
	return `f`
}

// UpsertSQL implements qb.Driver
func (d Driver) UpsertSQL(_ *qb.Table, _ []qb.Field, _ qb.Query) (string, []interface{}) {
	panic(`This should not be used`)
}

// Limit implements qb.Driver
func (d Driver) Limit(sql qb.SQL, limit int) {
	sql.WriteLine(`LIMIT ` + strconv.Itoa(limit))
}

// Returning implements qb.Driver
func (d Driver) Returning(b qb.SQLBuilder, q qb.Query, f []qb.Field) (string, []interface{}) {
	panic(`This should not be used`)
}

var types = map[qb.DataType]string{
	qb.Int:    `int`,
	qb.String: `string`,
	qb.Bool:   `boolean`,
	qb.Float:  `float`,
	qb.Date:   `date`,
	qb.Time:   `time`,
}

// TypeName returns the sql name for a type
func (d Driver) TypeName(t qb.DataType) string {
	if s, ok := types[t]; ok {
		return s
	}
	panic(`Unknown type`)
}

// Override returns the override map
func (d Driver) Override() qb.OverrideMap {
	return qb.OverrideMap{}
}
