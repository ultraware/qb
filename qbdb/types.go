package qbdb

import (
	"database/sql"

	"git.ultraware.nl/NiseVoid/qb"
)

// Target is a target for a query, either a plain DB or a Tx
type Target interface {
	Render(qb.Query) (string, []interface{})
	Query(qb.SelectQuery) (Rows, error)
	MustQuery(qb.SelectQuery) Rows
	QueryRow(qb.SelectQuery) Row
	Exec(q qb.Query) (Result, error)
	MustExec(q qb.Query) Result
	RawExec(s string, v ...interface{}) (Result, error)
	MustRawExec(s string, v ...interface{}) Result
	Prepare(q qb.Query) (*Stmt, error)
	MustPrepare(q qb.Query) *Stmt
}

// Tx is a transaction
type Tx struct {
	QueryTarget
	tx *sql.Tx
}

// Commit applies all the changes from the transaction
func (t Tx) Commit() error {
	return t.tx.Commit()
}

// MustCommit is the same as Commit, but it panics if an error occurred
func (t Tx) MustCommit() {
	err := t.Commit()
	if err != nil {
		panic(err)
	}
}

// Rollback reverts all the changes from the transaction
func (t Tx) Rollback() error {
	return t.tx.Rollback()
}

// QueryTarget is a sql.DB or sql.Tx
type QueryTarget struct {
	src interface {
		Exec(string, ...interface{}) (sql.Result, error)
		Query(string, ...interface{}) (*sql.Rows, error)
		QueryRow(string, ...interface{}) *sql.Row
		Prepare(string) (*sql.Stmt, error)
	}
	Driver qb.Driver
	Debug  bool
}

// DB wraps sql.DB to support qb queries
type DB struct {
	QueryTarget
	DB *sql.DB
}

// Begin starts a transaction
func (db *DB) Begin() (Tx, error) {
	tx, err := db.DB.Begin()
	return Tx{QueryTarget{tx, db.QueryTarget.Driver, db.QueryTarget.Debug}, tx}, err
}

// MustBegin is the same as Begin, but it panics if an error occurred
func (db *DB) MustBegin() Tx {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}

// New returns a new DB
func New(driver qb.Driver, db *sql.DB) *DB {
	return &DB{QueryTarget{db, driver, false}, db}
}

///////// Wrappers for Must functions /////////

// Rows is a wrapper for sql.Rows that adds MustScan
type Rows struct {
	*sql.Rows
}

// MustScan is the same as Scan except if an error occurs returned it will panic
func (r Rows) MustScan(dest ...interface{}) {
	err := r.Scan(dest...)
	if err != nil {
		panic(err)
	}
}

// Row is a wrapper for sql.Row that adds MustScan
type Row struct {
	*sql.Row
}

// MustScan returns true if there was a row.
/// If an error occurs it will panic
func (r Row) MustScan(dest ...interface{}) bool {
	err := r.Scan(dest...)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return err == nil
}

// Result is a wrapper for sql.Result that adds MustLastInsertId and MustRowsAffected
type Result struct {
	sql.Result
}

// MustLastInsertId is the same as LastInsertId except if an error occurs returned it will panic
func (r Result) MustLastInsertId() int64 { //nolint: stylecheck, golint
	id, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}

// MustRowsAffected is the same as RowsAffected except if an error occurs returned it will panic
func (r Result) MustRowsAffected() int64 {
	rows, err := r.RowsAffected()
	if err != nil {
		panic(err)
	}
	return rows
}
