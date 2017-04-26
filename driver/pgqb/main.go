package pgqb

import (
	"database/sql"
	"strconv"
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
)

func printType(v interface{}, c *int) (string, bool) {
	switch t := v.(type) {
	case int, int8, int16, int32, int64:
		return printInt(t), false
	case uint, uint8, uint16, uint32, uint64:
		return printUint(t), false
	case float32:
		return strconv.FormatFloat(float64(t), 'g', -1, 64), false
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64), false
	case bool:
		return strconv.FormatBool(t), false
	default:
		(*c)++
		return `$` + strconv.Itoa(*c), true
	}
}

func printInt(i interface{}) string {
	var v int64
	switch t := i.(type) {
	case int:
		v = int64(t)
	case int8:
		v = int64(t)
	case int16:
		v = int64(t)
	case int32:
		v = int64(t)
	case int64:
		v = t
	}
	return strconv.FormatInt(v, 10)
}

func printUint(i interface{}) string {
	var v uint64
	switch t := i.(type) {
	case uint:
		v = uint64(t)
	case uint8:
		v = uint64(t)
	case uint16:
		v = uint64(t)
	case uint32:
		v = uint64(t)
	case uint64:
		v = t
	}
	return strconv.FormatUint(v, 10)
}

func prepareQuery(q qb.SelectQuery) (string, []interface{}) {
	s, v := q.SQL()

	return prepareSQL(s, v)
}

func prepareSQL(s string, v []interface{}) (string, []interface{}) {
	vc := 0
	c := 0

	new := []interface{}{}
	for {
		i := strings.IndexRune(s, '?')
		if i == -1 {
			break
		}
		vc++
		str, param := printType(v[vc-1], &c)
		if param {
			new = append(new, v[vc-1])
		}
		s = s[:i] + str + s[i+1:]
	}

	return s, new
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
func (db *DB) Query(q qb.SelectQuery) (*qb.Cursor, error) {
	s, v := prepareQuery(q)

	r, err := db.DB.Query(s, v...)
	return qb.NewCursor(q.Fields(), r), err
}

// QueryRow executes the givven SelectQuery on the database, only returns one row
func (db *DB) QueryRow(q qb.SelectQuery) error {
	s, v := prepareQuery(q)
	r := db.DB.QueryRow(s, v...)

	f := q.Fields()
	dst := make([]interface{}, len(f))
	for k, v := range f {
		dst[k] = v
	}
	return r.Scan(dst...)
}

// RawExec executes the given SQL with the given params directly on the database
func (db *DB) RawExec(s string, v ...interface{}) error {
	_, err := db.DB.Exec(s, v...)
	return err
}

type savable interface {
	InsertHeader([]qb.DataField) string
	InsertValues([]qb.DataField) (string, []interface{})
	UpdateRecord([]qb.DataField) (string, []interface{})
	All() []qb.DataField
}

// Update updates the record in the database
func (db *DB) Update(record savable) error {
	s, v := prepareSQL(record.UpdateRecord(record.All()))
	_, err := db.DB.Exec(s, v...)
	return err
}

// Insert inster the record in the database
func (db *DB) Insert(record savable) error {
	s, v := record.InsertValues(record.All())
	s, v = prepareSQL(record.InsertHeader(record.All())+s, v)

	_, err := db.DB.Exec(s, v...)
	return err
}

// BatchInsert is used to insert multiple records at once
type BatchInsert struct {
	count     int
	SQL       string
	Values    []interface{}
	conflict  []qb.DataField
	updatable []qb.DataField
}

// NewBatch returns a BatchInsert
func (db *DB) NewBatch(record savable, conflict ...qb.DataField) *BatchInsert {
	s := record.InsertHeader(record.All())
	return &BatchInsert{0, s, nil, conflict, nil}
}

// Add adds the given record to the batch
func (b *BatchInsert) Add(record savable) {
	updatable := qb.GetUpdatableFields(record.All())
	if len(updatable) > len(b.updatable) {
		b.updatable = updatable
	}

	s, v := record.InsertValues(record.All())
	if b.count > 0 {
		s = `, ` + s
	}
	b.SQL += s
	b.Values = append(b.Values, v...)
	b.count++
}

// ExecuteBatch executes the given batch
func (db *DB) ExecuteBatch(b *BatchInsert) error {
	if b.count == 0 {
		panic(`Cannot insert empty batch`)
	}
	s, v := b.SQL, b.Values
	if len(b.conflict) > 0 {
		s += qb.GetUpsertSQL(b.conflict, b.updatable)
	}
	s, v = prepareSQL(s, v)

	_, err := db.DB.Exec(s, v...)
	return err
}

// Upsert tries to insert a record or update if a given field conflicts
func (db *DB) Upsert(record savable, conflict ...qb.DataField) error {
	s, v := record.InsertValues(record.All())
	s = record.InsertHeader(record.All()) + s
	s += qb.GetUpsertSQL(conflict, record.All())
	s, v = prepareSQL(s, v)
	_, err := db.DB.Exec(s, v...)
	return err
}
