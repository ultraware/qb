package msqb // import "git.ultraware.nl/Ultraware/qb/driver/msqb"

import (
	"database/sql"
	"strconv"
	"strings"

	"git.ultraware.nl/Ultraware/qb"
	"git.ultraware.nl/Ultraware/qb/driver/msqb/msqf"
	"git.ultraware.nl/Ultraware/qb/qbdb"
	"git.ultraware.nl/Ultraware/qb/qf"
)

// Driver implements MSSQL-specific features
type Driver struct{}

// New returns the driver
func New(db *sql.DB) qbdb.DB {
	return qbdb.New(Driver{}, db)
}

// ValueString returns a the SQL for a parameter value
func (d Driver) ValueString(i int) string {
	return `@p` + strconv.Itoa(i)
}

// BoolString formats a boolean in a format supported by MSSQL
func (d Driver) BoolString(v bool) string {
	if v {
		return `1`
	}
	return `0`
}

// EscapeCharacter returns the correct escape character for MSSQL
func (d Driver) EscapeCharacter() string {
	return `"`
}

// UpsertSQL implements qb.Driver
func (d Driver) UpsertSQL(t *qb.Table, _ []qb.Field, q qb.Query) (string, []interface{}) {
	panic(`mssql does not support upsert`)
}

// IgnoreConflictSQL implements qb.Driver
func (d Driver) IgnoreConflictSQL(_ *qb.Table, _ []qb.Field) (string, []interface{}) {
	panic(`mssql does not support ignore conflicts`)
}

// LimitOffset implements qb.Driver
func (d Driver) LimitOffset(sql qb.SQL, limit, offset int) {
	if offset > 0 {
		sql.WriteLine(`OFFSET ` + strconv.Itoa(offset) + ` ROWS`)
		if limit > 0 {
			sql.WriteLine(`FETCH NEXT ` + strconv.Itoa(limit) + ` ROWS ONLY`)
		}
		return
	}

	if limit > 0 {
		s := sql.String()

		sql.Rewrite(`SELECT TOP ` + strconv.Itoa(limit) + s[6:])
	}
}

// Returning implements qb.Driver
func (d Driver) Returning(b qb.SQLBuilder, q qb.Query, f []qb.Field) (string, []interface{}) {
	sql, v := q.SQL(b)

	t, insertBefore := `INSERTED`, `FROM`

	switch strings.SplitN(sql, ` `, 2)[0] {
	case `DELETE`:
		t = `DELETED`
	case `INSERT`:
		insertBefore = `VALUES`
	}

	line := ``
	for k, field := range f {
		if k > 0 {
			line += `, `
		}
		line += t + `.` + field.(*qb.TableField).Name
	}

	index := strings.Index(sql, insertBefore)
	if index < 0 {
		sql = sql + `OUTPUT ` + line + qb.NEWLINE
	} else {
		sql = sql[:index] + `OUTPUT ` + line + qb.NEWLINE + sql[index:]
	}

	return sql, v
}

var types = map[qb.DataType]string{
	qb.Int:    `int`,
	qb.String: `text`,
	qb.Bool:   `bit`,
	qb.Float:  `float`,
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
	override.Add(qf.Concat, msqf.Concat)
	override.Add(qf.Extract, msqf.DatePart)
	override.Add(qf.Now, msqf.GetDate)
}

// Override implements qb.Driver
func (d Driver) Override() qb.OverrideMap {
	return override
}
