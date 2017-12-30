package qbdb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
)

func (db QueryTarget) printType(v interface{}, c *int) (string, bool) {
	if _, ok := v.(driver.Valuer); ok {
		(*c)++
		return db.Driver.ValueString(*c), true
	}
	if v == nil {
		return `NULL`, false
	}
	switch t := reflect.ValueOf(v); t.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(t.Int(), 10), false
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(t.Uint(), 10), false
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(t.Float(), 'g', -1, 64), false
	case reflect.Bool:
		return db.Driver.BoolString(t.Bool()), false
	default:
		(*c)++
		return db.Driver.ValueString(*c), true
	}
}

// Render returns the generated SQL and values without executing the query
func (db QueryTarget) Render(q qb.Query) (string, []interface{}) {
	return db.prepare(q)
}

func (db QueryTarget) ctes(ctes []*qb.CTE, done map[*qb.CTE]bool, b qb.SQLBuilder) []string {
	var list []string

	for _, v := range ctes {
		if _, ok := done[v]; ok {
			continue
		}
		done[v] = true

		tmp := b.Context.Values

		newValues, newCTEs := []interface{}{}, []*qb.CTE{}
		b.Context.Values, b.Context.CTEs = &newValues, &newCTEs

		new := v.With(b)

		b.Context.Values = tmp
		list = append(append(list, db.ctes(*b.Context.CTEs, done, b)...), new)

		b.Context.Add(newValues...)
	}

	return list
}

func (db QueryTarget) prepare(q qb.Query) (string, []interface{}) {
	b := qb.NewSQLBuilder(db.Driver)

	s, v := q.SQL(b)

	if _, ok := q.(qb.SelectQuery); ok {
		ctes := b.Context.CTEs

		newList := []interface{}{}
		b.Context.Values = &newList
		if len(*ctes) > 0 {
			done := make(map[*qb.CTE]bool)

			s = `WITH ` + strings.Join(db.ctes(*ctes, done, b), `, `) + "\n\n" + s
		}
		v = append(*b.Context.Values, v...)
	}

	return db.prepareSQL(s, v)
}

func (db QueryTarget) prepareSQL(s string, v []interface{}) (string, []interface{}) {
	vc := 0
	c := 0
	offset := 0

	new := []interface{}{}
	for {
		i := strings.IndexRune(s[offset:], '?')
		if i == -1 {
			break
		}

		offset += i

		vc++
		str, param := db.printType(v[vc-1], &c)
		if param {
			new = append(new, v[vc-1])
		}

		if str != `?` {
			s = s[:offset] + str + s[offset+1:]
			offset += len(str) - 1
		}
		offset++
	}

	db.log(s, new)
	return s, new
}

func (db QueryTarget) log(s string, v []interface{}) {
	if db.Debug {
		fmt.Printf("-- Running query:\n%s-- With values: %v\n\n", s, v)
	}
}

// Query executes the given SelectQuery on the database
func (db QueryTarget) Query(q qb.SelectQuery) (*sql.Rows, error) {
	s, v := db.prepare(q)
	return db.src.Query(s, v...)
}

// QueryRow executes the given SelectQuery on the database, only returns one row
func (db QueryTarget) QueryRow(q qb.SelectQuery) *sql.Row {
	s, v := db.prepare(q)
	return db.src.QueryRow(s, v...)
}

// Exec executes the given query, returns only an error
func (db QueryTarget) Exec(q qb.Query) error {
	s, v := db.prepare(q)
	return db.RawExec(s, v...)
}

// RawExec executes the given SQL with the given params directly on the database
func (db QueryTarget) RawExec(s string, v ...interface{}) error {
	_, err := db.src.Exec(s, v...)
	return err
}

// Prepare prepares a query for efficient repeated executions
func (db QueryTarget) Prepare(q qb.Query) (*Stmt, error) {
	s, v := db.prepare(q)

	stmt, err := db.src.Prepare(s)
	if err != nil {
		return nil, err
	}

	return &Stmt{stmt, v}, nil
}
