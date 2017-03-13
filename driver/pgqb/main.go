package pgqb

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"git.ultraware.nl/NiseVoid/qb"
)

func getType(v interface{}) string {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return `::int`
	case float32, float64:
		return `::float`
	case string:
		return `::text`
	case bool:
		return `::bool`
	case time.Time:
		return `::timestamp`
	default:
		return ``
	}
}

func prepareQuery(q qb.SelectQuery) (string, []interface{}) {
	s, v := q.SQL()

	c := 0
	for {
		i := strings.IndexRune(s, '?')
		if i == -1 {
			break
		}
		c++
		s = s[:i] + `$` + strconv.Itoa(c) + getType(v[c-1]) + s[i+1:]
	}

	return s, v
}

// DB ...
type DB struct {
	DB *sql.DB
}

// New returns a new postgresql qb wrapper
func New(db *sql.DB) *DB {
	return &DB{DB: db}
}

// Query executes the givven SelectQuery on the database
func (db *DB) Query(q qb.SelectQuery) (*sql.Rows, error) {
	s, v := prepareQuery(q)
	return db.DB.Query(s, v...)
}

// QueryRow executes the givven SelectQuery on the database, only returns one row
func (db *DB) QueryRow(q qb.SelectQuery) *sql.Row {
	s, v := prepareQuery(q)
	return db.DB.QueryRow(s, v...)
}

// RawExec executes the given SQL with the given params directly on the database
func (db *DB) RawExec(s string, v ...interface{}) error {
	_, err := db.DB.Exec(s, v...)
	return err
}
