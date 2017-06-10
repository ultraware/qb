package qbdb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
)

func (db QueryTarget) printType(v interface{}, c *int) (string, bool) {
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

func (db QueryTarget) prepare(q qb.Query) (string, []interface{}) {
	s, v := q.SQL(db.Driver)

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
func (db QueryTarget) Query(q qb.SelectQuery) (*qb.Cursor, error) {
	s, v := db.prepare(q)
	r, err := db.src.Query(s, v...)

	return qb.NewCursor(q.Fields(), r), err
}

// QueryRow executes the given SelectQuery on the database, only returns one row
func (db QueryTarget) QueryRow(q qb.SelectQuery) error {
	s, v := db.prepare(q)
	r := db.src.QueryRow(s, v...)

	return qb.ScanToFields(q.Fields(), r)
}

// Exec executes the given query, returns only an error
func (db QueryTarget) Exec(q qb.Query) error {
	s, v := db.prepare(q)
	_, err := db.src.Exec(s, v...)

	return err
}

// RawExec executes the given SQL with the given params directly on the database
func (db QueryTarget) RawExec(s string, v ...interface{}) error {
	_, err := db.src.Exec(s, v...)
	return err
}
