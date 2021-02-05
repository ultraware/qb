package myqb // import "git.ultraware.nl/NiseVoid/qb/driver/myqb"

import (
	"database/sql"
	"strconv"
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/myqb/myqf"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qf"
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

// BoolString formats a boolean in a format supported by MySQL
func (d Driver) BoolString(v bool) string {
	if v {
		return `1`
	}
	return `0`
}

// EscapeCharacter returns the correct escape character for MySQL
func (d Driver) EscapeCharacter() string {
	return "`"
}

// UpsertSQL implements qb.Driver
func (d Driver) UpsertSQL(t *qb.Table, _ []qb.Field, q qb.Query) (string, []interface{}) {
	usql, values := q.SQL(qb.NewSQLBuilder(d))
	if !strings.HasPrefix(usql, `UPDATE `+t.Name) {
		panic(`Update does not update the correct table`)
	}
	usql = strings.ReplaceAll(usql, `UPDATE `+t.Name+qb.NEWLINE+`SET`, `UPDATE`)

	return `ON DUPLICATE KEY ` + usql, values
}

// IgnoreConflictSQL implements qb.Driver
func (d Driver) IgnoreConflictSQL(_ *qb.Table, _ []qb.Field) (string, []interface{}) {
	panic(`mysql does not support ignore conflicts`)
}

// LimitOffset implements qb.Driver
func (d Driver) LimitOffset(sql qb.SQL, limit, offset int) { //nolint: dupl
	if limit > 0 {
		sql.WriteLine(`LIMIT ` + strconv.Itoa(limit))
	}
	if offset > 0 {
		sql.WriteLine(`OFFSET ` + strconv.Itoa(offset))
	}
}

// Returning implements qb.Driver
func (d Driver) Returning(b qb.SQLBuilder, q qb.Query, f []qb.Field) (string, []interface{}) {
	panic(`mysql does not support RETURNING`)
}

var types = map[qb.DataType]string{
	qb.Int:    `int`,
	qb.String: `text`,
	qb.Bool:   `boolean`,
	qb.Float:  `double`,
	qb.Date:   `date`,
	qb.Time:   `datetime`,
}

// TypeName implements qb.Driver
func (d Driver) TypeName(t qb.DataType) string {
	if s, ok := types[t]; ok {
		return s
	}
	panic(`Unknown type`)
}

var override = qb.OverrideMap{}

func init() {
	override.Add(qf.Excluded, myqf.Values)
}

// Override implements qb.Driver
func (d Driver) Override() qb.OverrideMap {
	return override
}
