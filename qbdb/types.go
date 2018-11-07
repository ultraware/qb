package qbdb

import (
	"database/sql"

	"github.com/ultraware/qb"
)

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
