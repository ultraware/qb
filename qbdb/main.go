package qbdb // import "git.ultraware.nl/Ultraware/qb/v2/qbdb"

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"git.ultraware.nl/Ultraware/qb/v2"
)

func (db queryTarget) printType(v interface{}, c *int) (string, bool) {
	if _, ok := v.(driver.Valuer); ok {
		(*c)++
		return db.driver.ValueString(*c), true
	}
	if v == nil {
		return `NULL`, false
	}

	switch t := reflect.ValueOf(v); t.Type().Kind() { //nolint: exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(t.Int(), 10), false
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(t.Uint(), 10), false
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(t.Float(), 'g', -1, 64), false
	case reflect.Bool:
		return db.driver.BoolString(t.Bool()), false
	default:
		(*c)++
		return db.driver.ValueString(*c), true
	}
}

// Render returns the generated SQL and values without executing the query
func (db queryTarget) Render(q qb.Query) (string, []interface{}) {
	return db.prepare(q)
}

func (db queryTarget) ctes(ctes []*qb.CTE, done map[*qb.CTE]bool, b qb.SQLBuilder) []string {
	var list []string

	for _, v := range ctes {
		if _, ok := done[v]; ok {
			continue
		}
		done[v] = true

		tmp := b.Context.Values

		newValues, newCTEs := []interface{}{}, []*qb.CTE{}
		b.Context.Values, b.Context.CTEs = &newValues, &newCTEs

		newWith := v.With(b)

		b.Context.Values = tmp
		list = append(append(list, db.ctes(*b.Context.CTEs, done, b)...), newWith)

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

func (db queryTarget) prepare(q qb.Query) (string, []interface{}) {
	b := qb.NewSQLBuilder(db.driver)

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

func (db queryTarget) prepareSQL(s string, v []interface{}) (string, []interface{}) {
	vc := 0
	c := 0
	b := &strings.Builder{}

	newValue := []interface{}{}
	for _, chr := range s {
		if chr != '?' {
			if _, err := b.WriteRune(chr); err != nil {
				panic(err)
			}
			continue
		}

		str, param := db.printType(v[vc], &c)
		if param {
			newValue = append(newValue, v[vc])
		}
		vc++

		if _, err := b.WriteString(str); err != nil {
			panic(err)
		}
	}

	db.log(b.String(), newValue)
	return b.String(), newValue
}

func (db queryTarget) log(s string, v []interface{}) {
	if db.debug {
		fmt.Printf("-- Running query:\n%s-- With values: %v\n\n", s, v)
	}
}

// Query executes the given SelectQuery on the database
func (db queryTarget) Query(q qb.SelectQuery) (Rows, error) {
	return db.QueryContext(context.Background(), q)
}

// QueryContext executes the given SelectQuery on the database
func (db queryTarget) QueryContext(c context.Context, q qb.SelectQuery) (Rows, error) {
	s, v := db.prepare(q)
	r, err := db.RawQueryContext(c, s, v...)
	return r, err
}

// MustQuery executes the given SelectQuery on the database
// If an error occurs returned it will panic
func (db queryTarget) MustQuery(q qb.SelectQuery) Rows {
	return db.MustQueryContext(context.Background(), q)
}

// MustQueryContext executes the given SelectQuery on the database
// If an error occurs returned it will panic
func (db queryTarget) MustQueryContext(c context.Context, q qb.SelectQuery) Rows {
	return mustRows(db.QueryContext(c, q))
}

// RawQuery executes the given raw query on the database
func (db queryTarget) RawQuery(s string, v ...interface{}) (Rows, error) {
	return db.RawQueryContext(context.Background(), s, v...)
}

// RawQueryContext executes the given raw query on the database
func (db queryTarget) RawQueryContext(c context.Context, s string, v ...interface{}) (Rows, error) {
	r, err := db.src.QueryContext(c, s, v...)
	return Rows{r}, err
}

// MustRawQuery executes the given raw query on the database
// If an error occurs returned it will panic
func (db queryTarget) MustRawQuery(s string, v ...interface{}) Rows {
	return db.MustRawQueryContext(context.Background(), s, v...)
}

// MustRawQueryContext executes the given raw query on the database
// If an error occurs returned it will panic
func (db queryTarget) MustRawQueryContext(c context.Context, s string, v ...interface{}) Rows {
	return mustRows(db.RawQueryContext(c, s, v...))
}

// QueryRow executes the given SelectQuery on the database, only returns one row
func (db queryTarget) QueryRow(q qb.SelectQuery) Row {
	return db.QueryRowContext(context.Background(), q)
}

// QueryRowContext executes the given SelectQuery on the database, only returns one row
func (db queryTarget) QueryRowContext(c context.Context, q qb.SelectQuery) Row {
	if sq, ok := q.(*qb.SelectBuilder); ok {
		sq.Limit(1)
	}

	s, v := db.prepare(q)
	return db.RawQueryRowContext(c, s, v...)
}

// RawQueryRow executes the given raw query on the database, only returns one row
func (db queryTarget) RawQueryRow(s string, v ...interface{}) Row {
	return db.RawQueryRowContext(context.Background(), s, v...)
}

// RawQueryRowContext executes the given raw query on the database, only returns one row
func (db queryTarget) RawQueryRowContext(c context.Context, s string, v ...interface{}) Row {
	return Row{db.src.QueryRowContext(c, s, v...)}
}

// Exec executes the given query, returns only an error
func (db queryTarget) Exec(q qb.Query) (Result, error) {
	return db.ExecContext(context.Background(), q)
}

// ExecContext executes the given query, returns only an error
func (db queryTarget) ExecContext(c context.Context, q qb.Query) (Result, error) {
	s, v := db.prepare(q)
	return db.RawExecContext(c, s, v...)
}

// MustExec executes the given query
// If an error occurs returned it will panic
func (db queryTarget) MustExec(q qb.Query) Result {
	return db.MustExecContext(context.Background(), q)
}

// MustExecContext executes the given query with context
// If an error occurs returned it will panic
func (db queryTarget) MustExecContext(c context.Context, q qb.Query) Result {
	return mustResult(db.ExecContext(c, q))
}

// RawExec executes the given SQL with the given params directly on the database
func (db queryTarget) RawExec(s string, v ...interface{}) (Result, error) {
	return db.RawExecContext(context.Background(), s, v...)
}

// RawExecContext executes the given SQL with the given params directly on the database
func (db queryTarget) RawExecContext(c context.Context, s string, v ...interface{}) (Result, error) {
	r, err := db.src.ExecContext(c, s, v...)
	return Result{r}, err
}

// MustRawExec executes the given SQL with the given params directly on the database
// If an error occurs returned it will panic
func (db queryTarget) MustRawExec(s string, v ...interface{}) Result {
	return db.MustRawExecContext(context.Background(), s, v...)
}

// MustRawExecContext executes the given SQL with the given params directly on the database
// If an error occurs returned it will panic
func (db queryTarget) MustRawExecContext(c context.Context, s string, v ...interface{}) Result {
	return mustResult(db.RawExecContext(c, s, v...))
}

// Prepare prepares a query for efficient repeated executions
func (db queryTarget) Prepare(q qb.Query) (*Stmt, error) {
	return db.PrepareContext(context.Background(), q)
}

// PrepareContext prepares a query for efficient repeated executions
func (db queryTarget) PrepareContext(c context.Context, q qb.Query) (*Stmt, error) {
	s, v := db.prepare(q)

	stmt, err := db.src.PrepareContext(c, s)
	if err != nil {
		return nil, err
	}

	return &Stmt{stmt, v}, nil
}

// MustPrepare prepares a query for efficient repeated executions
// If an error occurs returned it will panic
func (db queryTarget) MustPrepare(q qb.Query) *Stmt {
	return db.MustPrepareContext(context.Background(), q)
}

// MustPrepareContext prepares a query for efficient repeated executions
// If an error occurs returned it will panic
func (db queryTarget) MustPrepareContext(ctx context.Context, q qb.Query) *Stmt {
	stmt, err := db.PrepareContext(ctx, q)
	if err != nil {
		panic(err)
	}
	return stmt
}

func mustRows(r Rows, err error) Rows {
	if err != nil {
		panic(err)
	}
	return r
}

func mustResult(r Result, err error) Result {
	if err != nil {
		panic(err)
	}
	return r
}
