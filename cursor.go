package qb

type rows interface {
	Next() bool
	Scan(...interface{}) error
	Close() error
}

// Cursor ...
type Cursor struct {
	fields             []DataField
	rows               rows
	DisableExitOnError bool
	err                error
}

// NewCursor returns a new Cursor
func NewCursor(f []DataField, r rows) *Cursor {
	c := &Cursor{fields: f, rows: r}
	return c
}

// Next loads the next row into the fields
func (c *Cursor) Next() bool {
	if !c.rows.Next() {
		c.Close()
		return false
	}
	err := ScanToFields(c.fields, c.rows)
	if err != nil {
		c.err = err
		if !c.DisableExitOnError {
			c.Close()
		}
		return false
	}

	return true
}

// Close the rows object, automatically called by Next when all rows have been read
func (c *Cursor) Close() {
	err := c.rows.Close()
	if c.err == nil {
		c.err = err
	}
}

// Error returns the last error
func (c Cursor) Error() error {
	return c.err
}
