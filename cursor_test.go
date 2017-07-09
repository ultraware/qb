package qb

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRows struct {
	Closed   bool
	Position int
	Data     [][]interface{}
	Err      error
}

func newRows(data [][]interface{}) *testRows {
	return &testRows{false, -1, data, nil}
}

func (r *testRows) Next() bool {
	if r.Position >= len(r.Data)-1 {
		return false
	}
	r.Position++
	return true
}

func (r *testRows) Scan(dst ...interface{}) error {
	if r.Position > len(r.Data)-1 {
		return errors.New(`All testRows passed`)
	}

	if len(dst) != len(r.Data[r.Position]) {
		return errors.New(`Invalid number of arguments`)
	}

	for k, v := range dst {
		str := r.Data[r.Position][k].(string)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(&str))
	}
	return nil
}

func (r *testRows) Close() error {
	r.Closed = true
	r.Data = nil
	return r.Err
}

func TestCursor(t *testing.T) {
	assert := assert.New(t)

	s1, s2 := ``, ``
	f1, f2 := NewDataField(nil, &s1), NewDataField(nil, &s2)

	rows := newRows([][]interface{}{
		{`0.0`, `1.0`},
		{`0.1`, `1.1`},
		{`0.2`, `1.2`},
	})
	cur := NewCursor([]DataField{f1, f2}, rows)

	assert.True(cur.Next())
	assert.True(`0.0` == s1 && `1.0` == s2)
	assert.True(cur.Next())
	assert.True(`0.1` == s1 && `1.1` == s2)
	assert.True(cur.Next())
	assert.True(`0.2` == s1 && `1.2` == s2)

	assert.False(cur.Next())

	cur.Close()
	assert.NoError(cur.Error())
	assert.True(rows.Closed)
}

func TestCloseErr(t *testing.T) {
	assert := assert.New(t)

	rows := newRows(nil)
	cur := NewCursor(nil, rows)

	rows.Err = errors.New(``)

	cur.Close()
	assert.Error(cur.Error())
}

func TestScanErr(t *testing.T) {
	assert := assert.New(t)

	rows := newRows([][]interface{}{
		{`0.0`, `1.0`},
		{`0.1`, `1.1`},
		{`0.2`, `1.2`},
	})
	cur := NewCursor(nil, rows)

	cur.DisableExitOnError = true
	cur.Next()
	assert.Error(cur.Error())
	assert.False(rows.Closed)

	cur.DisableExitOnError = false
	cur.Next()
	assert.True(rows.Closed)
}
