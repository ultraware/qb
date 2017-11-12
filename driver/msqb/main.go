package msqb

import (
	"database/sql"
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

// UpsertSQL ...
func (d Driver) UpsertSQL(t *qb.Table, _ []qb.Field, q qb.Query) (string, []interface{}) {
	panic(`mssql does not support upsert`)
}

// Returning ...
func (d Driver) Returning(q qb.Query, f []qb.Field) (string, []interface{}) {
	c := qb.NewContext(d, qb.NoAlias())
	sql, v := q.SQL(d)

	var t string
	switch strings.SplitN(sql, ` `, 2)[0] {
	case `DELETE`:
		t = `DELETED`
	default:
		t = `INSERTED`
	}

	line := ``
	for k, field := range f {
		if k > 0 {
			line += `, `
		}
		line += t + `.` + field.QueryString(c)
	}

	index := strings.Index(sql, `WHERE`)
	if index < 0 {
		index = len(sql)
	}
	sql = sql[:index] + `OUTPUT ` + line + qb.NEWLINE + sql[index:]

	return sql, v
}

var types = map[qb.DataType]string{
	qb.Int:     `int`,
	qb.String:  `text`,
	qb.Boolean: `bit`,
	qb.Float:   `float`,
	qb.Date:    `date`,
	qb.Time:    `datetime`,
}

// TypeName ...
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

// Override ...
func (d Driver) Override() qb.OverrideMap {
	return override
}
