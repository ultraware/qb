package pgqb

import (
	"database/sql"
	"strconv"
	"strings"

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

// ConcatOperator ...
func (d Driver) ConcatOperator() string {
	return `||`
}

// ExcludedField ...
func (d Driver) ExcludedField(f string) string {
	return `EXCLUDED.` + f
}

// UpsertSQL ...
func (d Driver) UpsertSQL(t *qb.Table, conflict []qb.Field, q qb.Query) (string, []interface{}) {
	sql := ``
	for k, v := range conflict {
		if k > 0 {
			sql += qb.COMMA
		}
		sql += v.QueryString(d, qb.NoAlias(), nil)
	}

	usql, values := q.SQL(d)
	if !strings.HasPrefix(usql, `UPDATE `+t.Name) {
		panic(`Update does not update the correct table`)
	}
	usql = strings.Replace(usql, `UPDATE `+t.Name, `UPDATE`, -1)

	return `ON CONFLICT (` + sql + `) DO ` + usql, values
}

// Returning ...
func (d Driver) Returning(q qb.Query, f string) (string, []interface{}) {
	s, v := q.SQL(d)
	return s + `RETURNING ` + f, v
}

// DateExtract ...
func (d Driver) DateExtract(f string, part string) string {
	return `EXTRACT(` + part + ` FROM ` + f + `)`
}
