package msqb

import (
	"database/sql"
	"strconv"
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/msqb/msqf"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qf"
)

// Driver implements MSSQL-specific features
type Driver struct{}

// New returns the driver
func New(db *sql.DB) *qbdb.DB {
	return qbdb.New(Driver{}, db)
}

// ValueString returns a the SQL for a parameter value
func (d Driver) ValueString(i int) string {
	return `@p` + strconv.Itoa(i)
}

// BoolString formats a boolean in a format supported by PostgreSQL
func (d Driver) BoolString(v bool) string {
	if v {
		return `1`
	}
	return `0`
}

// UpsertSQL implements qb.Driver
func (d Driver) UpsertSQL(t *qb.Table, _ []qb.Field, q qb.Query) (string, []interface{}) {
	panic(`mssql does not support upsert`)
}

// Limit implements qb.Driver
func (d Driver) Limit(sql qb.SQL, limit int) {
	s := sql.String()

	sql.Rewrite(`SELECT TOP ` + strconv.Itoa(limit) + s[6:])
}

// Returning implements qb.Driver
func (d Driver) Returning(q qb.Query, f []qb.Field) (string, []interface{}) {
	b := qb.NewSQLBuilder(d)
	sql, v := q.SQL(b)

	t, insertBefore := `INSERTED`, `WHERE`

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
		line += t + `.` + field.QueryString(b.Context)
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
