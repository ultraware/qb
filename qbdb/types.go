package qbdb

import (
	"context"
	"database/sql"
	"errors"

	"git.ultraware.nl/Ultraware/qb/v2"
)

// Target is a target for a query, either a plain DB or a Tx
type Target interface {
	Render(qb.Query) (string, []interface{})
	Query(qb.SelectQuery) (Rows, error)
	QueryContext(ctx context.Context, q qb.SelectQuery) (Rows, error)
	RawQuery(string, ...interface{}) (Rows, error)
	RawQueryContext(ctx context.Context, s string, v ...interface{}) (Rows, error)
	MustQuery(qb.SelectQuery) Rows
	MustRawQuery(string, ...interface{}) Rows
	MustQueryContext(ctx context.Context, q qb.SelectQuery) Rows
	MustRawQueryContext(ctx context.Context, s string, v ...interface{}) Rows
	QueryRow(qb.SelectQuery) Row
	QueryRowContext(ctx context.Context, q qb.SelectQuery) Row
	RawQueryRow(string, ...interface{}) Row
	RawQueryRowContext(ctx context.Context, s string, v ...interface{}) Row
	Exec(q qb.Query) (Result, error)
	ExecContext(ctx context.Context, q qb.Query) (Result, error)
	RawExec(string, ...interface{}) (Result, error)
	RawExecContext(ctx context.Context, s string, v ...interface{}) (Result, error)
	MustExec(q qb.Query) Result
	MustExecContext(ctx context.Context, q qb.Query) Result
	MustRawExec(string, ...interface{}) Result
	MustRawExecContext(ctx context.Context, s string, v ...interface{}) Result
	Prepare(qb.Query) (*Stmt, error)
	PrepareContext(ctx context.Context, q qb.Query) (*Stmt, error)
	MustPrepare(qb.Query) *Stmt
	MustPrepareContext(ctx context.Context, q qb.Query) *Stmt
	Driver() qb.Driver
	SetDebug(bool)
}

// Tx is a transaction
type Tx interface {
	Target
	Commit() error
	MustCommit()
	Rollback() error
}

type tx struct {
	queryTarget
	tx *sql.Tx
}

// Commit applies all the changes from the transaction
func (t *tx) Commit() error {
	return t.tx.Commit()
}

// MustCommit is the same as Commit, but it panics if an error occurred
func (t *tx) MustCommit() {
	err := t.Commit()
	if err != nil {
		panic(err)
	}
}

// Rollback reverts all the changes from the transaction
func (t *tx) Rollback() error {
	return t.tx.Rollback()
}

type queryTarget struct {
	src interface {
		Exec(string, ...interface{}) (sql.Result, error)
		Query(string, ...interface{}) (*sql.Rows, error)
		QueryRow(string, ...interface{}) *sql.Row
		Prepare(string) (*sql.Stmt, error)
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
		PrepareContext(context.Context, string) (*sql.Stmt, error)
	}
	driver qb.Driver
	debug  bool
}

func (db *queryTarget) Driver() qb.Driver {
	return db.driver
}

func (db *queryTarget) SetDebug(b bool) {
	db.debug = b
}

// DB is a database connection
type DB interface {
	Target
	Begin() (Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	MustBegin() Tx
}

type db struct {
	queryTarget
	DB *sql.DB
}

// Begin starts a transaction
func (db *db) Begin() (Tx, error) {
	return db.BeginTx(context.Background(), nil)
}

// BeginTx starts a transaction
func (db *db) BeginTx(c context.Context, o *sql.TxOptions) (Tx, error) {
	rawTx, err := db.DB.BeginTx(c, o)
	return &tx{queryTarget{rawTx, db.queryTarget.driver, db.queryTarget.debug}, rawTx}, err
}

// MustBegin is the same as Begin, but it panics if an error occurred
func (db *db) MustBegin() Tx {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}

// New returns a new DB
func New(driver qb.Driver, database *sql.DB) DB {
	return &db{queryTarget{database, driver, false}, database}
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
// If an error occurs it will panic
func (r Row) MustScan(dest ...interface{}) bool {
	err := r.Scan(dest...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	return err == nil
}

// Result is a wrapper for sql.Result that adds MustLastInsertId and MustRowsAffected
type Result struct {
	sql.Result
}

// MustLastInsertId is the same as LastInsertId except if an error occurs returned it will panic
func (r Result) MustLastInsertId() int64 { //nolint: stylecheck
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
