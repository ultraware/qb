package pgqb // import "git.ultraware.nl/Ultraware/qb/v2/driver/pgqb"

import (
	"database/sql"
	"strconv"
	"strings"

	"git.ultraware.nl/Ultraware/qb/v2"
	"git.ultraware.nl/Ultraware/qb/v2/driver/pgqb/pgqc"
	"git.ultraware.nl/Ultraware/qb/v2/driver/pgqb/pgqf"
	"git.ultraware.nl/Ultraware/qb/v2/qbdb"
	"git.ultraware.nl/Ultraware/qb/v2/qc"
	"git.ultraware.nl/Ultraware/qb/v2/qf"
)

// Driver implements PostgreSQL-specific features
type Driver struct{}

// New returns the driver
func New(db *sql.DB) qbdb.DB {
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

// EscapeCharacter returns the correct escape character for PostgreSQL
func (d Driver) EscapeCharacter() string {
	return `"`
}

// UpsertSQL implements qb.Driver
func (d Driver) UpsertSQL(t *qb.Table, conflict []qb.Field, q qb.Query) (string, []interface{}) {
	c := qb.NewContext(d, qb.NoAlias())
	sql := ``
	for k, v := range conflict {
		if k > 0 {
			sql += qb.COMMA
		}
		sql += v.QueryString(c)
	}

	usql, values := q.SQL(qb.NewSQLBuilder(d))
	if !strings.HasPrefix(usql, `UPDATE `+t.Name) {
		panic(`Update does not update the correct table`)
	}
	usql = strings.ReplaceAll(usql, `UPDATE `+t.Name, `UPDATE`)

	return `ON CONFLICT (` + sql + `) DO ` + usql, values
}

// IgnoreConflictSQL implements qb.Driver
func (d Driver) IgnoreConflictSQL(t *qb.Table, conflict []qb.Field) (string, []interface{}) {
	c := qb.NewContext(d, qb.NoAlias())
	sql := ``
	for k, v := range conflict {
		if k > 0 {
			sql += qb.COMMA
		}
		sql += v.QueryString(c)
	}

	return `ON CONFLICT (` + sql + ") DO NOTHING\n", *c.Values
}

// LimitOffset implements qb.Driver
func (d Driver) LimitOffset(sql qb.SQL, limit, offset int) {
	if limit > 0 {
		sql.WriteLine(`LIMIT ` + strconv.Itoa(limit))
	}
	if offset > 0 {
		sql.WriteLine(`OFFSET ` + strconv.Itoa(offset))
	}
}

// Returning implements qb.Driver
func (d Driver) Returning(b qb.SQLBuilder, q qb.Query, f []qb.Field) (string, []interface{}) {
	s, v := q.SQL(b)

	line := ``
	for k, field := range f {
		if k > 0 {
			line += `, `
		}
		line += field.QueryString(b.Context)
	}

	return s + `RETURNING ` + line + qb.NEWLINE, append(v, *b.Context.Values...)
}

// LateralJoin implements qb.Driver
func (d Driver) LateralJoin(c *qb.Context, s *qb.SubQuery) string {
	return `LATERAL ` + s.TableString(c)
}

var types = map[qb.DataType]string{
	qb.Int:    `int`,
	qb.String: `text`,
	qb.Bool:   `boolean`,
	qb.Float:  `float`,
	qb.Date:   `date`,
	qb.Time:   `timestamptz`,
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
	override.Add(qf.Excluded, pgqf.Excluded)

	override.Add(qc.Like, pgqc.ILike)
}

// Override implements qb.Driver
func (d Driver) Override() qb.OverrideMap {
	return override
}
