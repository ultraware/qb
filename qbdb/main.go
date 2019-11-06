package qbdb // import "git.ultraware.nl/NiseVoid/qb/qbdb"

import (
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
	var list []string // nolint: prealloc

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

func shouldContext(q qb.Query) bool {
	switch t := q.(type) {
	case *qb.UpdateBuilder, qb.DeleteBuilder:
		return true
	case qb.ReturningBuilder:
		return shouldContext(t.Query)
	}

	return false
}

func (db QueryTarget) prepare(q qb.Query) (string, []interface{}) {
	b := qb.NewSQLBuilder(db.Driver)

	if shouldContext(q) {
		b.Context = qb.NewContext(b.Context.Driver, qb.AliasGenerator())
	}

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
	b := &strings.Builder{}

	new := []interface{}{}
	for _, chr := range s {
		if chr != '?' {
			if _, err := b.WriteRune(chr); err != nil {
				panic(err)
			}
			continue
		}

		str, param := db.printType(v[vc], &c)
		if param {
			new = append(new, v[vc])
		}
		vc++

		if _, err := b.WriteString(str); err != nil {
			panic(err)
		}
	}

	db.log(b.String(), new)
	return b.String(), new
}

func (db QueryTarget) log(s string, v []interface{}) {
	if db.Debug {
		fmt.Printf("-- Running query:\n%s-- With values: %v\n\n", s, v)
	}
}

// Query executes the given SelectQuery on the database
func (db QueryTarget) Query(q qb.SelectQuery) (Rows, error) {
	s, v := db.prepare(q)
	r, err := db.src.Query(s, v...)
	return Rows{r}, err
}

// MustQuery executes the given SelectQuery on the database
// If an error occurs returned it will panic
func (db QueryTarget) MustQuery(q qb.SelectQuery) Rows {
	r, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	return r
}

// QueryRow executes the given SelectQuery on the database, only returns one row
func (db QueryTarget) QueryRow(q qb.SelectQuery) Row {
	if sq, ok := q.(*qb.SelectBuilder); ok {
		sq.Limit(1)
	}

	s, v := db.prepare(q)
	return Row{db.src.QueryRow(s, v...)}
}

// Exec executes the given query, returns only an error
func (db QueryTarget) Exec(q qb.Query) (Result, error) {
	s, v := db.prepare(q)
	return db.RawExec(s, v...)
}

// MustExec executes the given query
// If an error occurs returned it will panic
func (db QueryTarget) MustExec(q qb.Query) Result {
	res, err := db.Exec(q)
	if err != nil {
		panic(err)
	}
	return res
}

// RawExec executes the given SQL with the given params directly on the database
func (db QueryTarget) RawExec(s string, v ...interface{}) (Result, error) {
	r, err := db.src.Exec(s, v...)
	return Result{r}, err
}

// MustRawExec executes the given SQL with the given params directly on the database
// If an error occurs returned it will panic
func (db QueryTarget) MustRawExec(s string, v ...interface{}) Result {
	res, err := db.RawExec(s, v...)
	if err != nil {
		panic(err)
	}
	return res
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

// MustPrepare prepares a query for efficient repeated executions
// If an error occurs returned it will panic
func (db QueryTarget) MustPrepare(q qb.Query) *Stmt {
	stmt, err := db.Prepare(q)
	if err != nil {
		panic(err)
	}
	return stmt
}
