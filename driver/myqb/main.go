package myqb

import (
	"database/sql"
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
)

// Driver implements PostgreSQL-specific features
type Driver struct{}

// New returns the driver
func New(db *sql.DB) *qbdb.DB {
	_, _ = db.Exec(`SET SESSION sql_mode = 'ANSI'`)
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

// ConcatOperator ...
func (d Driver) ConcatOperator() string {
	return `||`
}

// ExcludedField ...
func (d Driver) ExcludedField(f string) string {
	return `VALUES(` + f + `)`
}

// UpsertSQL ...
func (d Driver) UpsertSQL(t *qb.Table, _ []qb.Field, q qb.Query) (string, []interface{}) {
	usql, values := q.SQL(d)
	if !strings.HasPrefix(usql, `UPDATE `+t.Name) {
		panic(`Update does not update the correct table`)
	}
	usql = strings.Replace(usql, `UPDATE `+t.Name+qb.NEWLINE+`SET`, `UPDATE`, -1)

	return `ON DUPLICATE KEY ` + usql, values
}

// Returning ...
func (d Driver) Returning(q qb.Query, f string) (string, []interface{}) {
	panic(`mysql does not support RETURNING`)
}

// DateExtract ...
func (d Driver) DateExtract(f string, part string) string {
	return `EXTRACT(` + part + ` FROM ` + f + `)`
}
