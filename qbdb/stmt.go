package qbdb

import (
	"context"
	"database/sql"
)

// Stmt represents a prepared statement in the database
type Stmt struct {
	stmt *sql.Stmt
	args []interface{}
}

// Close closes the prepared statement
func (s *Stmt) Close() error {
	return s.stmt.Close()
}

// Query executes the prepared statement on the database
func (s *Stmt) Query() (Rows, error) {
	return s.QueryContext(context.Background())
}

// QueryContext executes the prepared statement on the database
func (s *Stmt) QueryContext(c context.Context) (Rows, error) {
	r, err := s.stmt.QueryContext(c, s.args...)
	return Rows{r}, err
}

// MustQuery executes the prepared statement on the database
func (s *Stmt) MustQuery() Rows {
	r, err := s.Query()
	if err != nil {
		panic(err)
	}
	return r
}

// QueryRow executes the prepared statement on the database, only returns one row
func (s *Stmt) QueryRow() Row {
	return s.QueryRowContext(context.Background())
}

// QueryRowContext executes the prepared statement on the database, only returns one row
func (s *Stmt) QueryRowContext(c context.Context) Row {
	return Row{s.stmt.QueryRowContext(c, s.args...)}
}

// Exec executes the prepared statement
func (s *Stmt) Exec() (Result, error) {
	return s.ExecContext(context.Background())
}

// ExecContext executes the prepared statement
func (s *Stmt) ExecContext(c context.Context) (Result, error) {
	r, err := s.stmt.ExecContext(c, s.args...)

	return Result{r}, err
}

// MustExec executes the given SelectQuery on the database
func (s *Stmt) MustExec() Result {
	r, err := s.Exec()
	if err != nil {
		panic(err)
	}
	return r
}
