package qbdb

import (
	"database/sql"

	"git.ultraware.nl/NiseVoid/qb"
)

// Driver implements databse-specific features
type Driver interface {
	ValueString(int) string
	BoolString(bool) string
	UpsertSQL(fields []qb.DataField, conflict []qb.DataField) string
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
	}
	Driver Driver
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

// New returns a new DB
func New(driver Driver, db *sql.DB) *DB {
	return &DB{QueryTarget{db, driver, false}, db}
}

// Savable has all methods needed to run Insert/Update/Delete
type Savable interface {
	All() []qb.DataField
	GetTable() *qb.Table
}
