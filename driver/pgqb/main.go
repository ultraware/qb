package pgqb

import (
	"database/sql"
	"strconv"

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
func (d Driver) ValueString(c int) string {
	return `$` + strconv.Itoa(c)
}

// BoolString formats a boolean in a format supported by PostgreSQL
func (d Driver) BoolString(v bool) string {
	return strconv.FormatBool(v)
}

// UpsertSQL returns the
func (d Driver) UpsertSQL(f []qb.DataField, conflict []qb.DataField) string {
	sql := ``
	for k, v := range conflict {
		if k > 0 {
			sql += qb.COMMA
		}
		sql += v.QueryString(&qb.NoAlias{}, nil)
	}
	return `ON CONFLICT (` + sql + `) DO ` + updateExcluded(f)
}

func updateExcluded(f []qb.DataField) string {
	s := `UPDATE SET `
	for k, v := range f {
		if k > 0 {
			s += `, `
		}
		sql := v.QueryString(&qb.NoAlias{}, nil)
		s += sql + ` = EXCLUDED.` + sql
		continue
	}
	return s + "\n"
}
