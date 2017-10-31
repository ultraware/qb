package qbdb

import (
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

// Query executes the given SelectQuery on the database
func (s *Stmt) Query() (*sql.Rows, error) {
	return s.stmt.Query(s.args...)
}

// QueryRow executes the given SelectQuery on the database, only returns one row
func (s *Stmt) QueryRow() *sql.Row {
	return s.stmt.QueryRow(s.args...)
}

// Exec executes the given query, returns only an error
func (s *Stmt) Exec() error {
	_, err := s.stmt.Exec(s.args...)

	return err
}
