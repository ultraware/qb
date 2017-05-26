package myqb

import (
	"database/sql"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
)

// Driver implements PostgreSQL-specific features
type Driver struct{}

// New returns the driver
func New(db *sql.DB) *qbdb.DB {
	return qbdb.New(Driver{}, db)
}

// ValueString returns a the SQL for a parameter value
func (d Driver) ValueString(_ int) string {
	return `?`
}

// BoolString formats a boolean in a format supported by PostgreSQL
func (d Driver) BoolString(v bool) string {
	if v {
		return `1`
	}
	return `0`
}

// UpsertSQL returns the
func (d Driver) UpsertSQL(f []qb.DataField, _ []qb.DataField) string {
	s := ``
	for k, v := range f {
		if k > 0 {
			s += `, `
		}
		sql := v.QueryString(&qb.NoAlias{}, nil)
		s += sql + ` = ` + `VALUES(` + sql + `)`
	}

	return `ON DUPLICATE KEY UPDATE SET ` + s + "\n"
}
