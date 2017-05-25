package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testField = &TableField{Name: `test`}

func ok(f func()) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	ok = true
	f()

	return
}

///// NewDataField /////

func TestNewDataField(t *testing.T) {
	assert := assert.New(t)

	var value int
	f := NewDataField(testField, &value)

	assert.Equal(f.Field, testField)
	assert.Equal(f.Value, &value)
}

func TestNewDataFieldPanics(t *testing.T) {
	var value int

	if ok(func() { _ = NewDataField(testField, value) }) {
		t.Error(`NewDataField should panic when the provided value is not a pointer`)
	}
}

func BenchmarkNewDataField(b *testing.B) {
	value := 0
	for i := 0; i < b.N; i++ {
		_ = NewDataField(testField, &value)
	}
}

///// Scan /////

func TestScanWithPointer(t *testing.T) {
	assert := assert.New(t)

	var value int
	f := NewDataField(testField, &value)

	test := 3

	v := f.getScanTarget()
	*v.(**int) = &test
	assert.NotEqual(test, value)

	f.updateData(v)

	assert.Equal(test, value)
	assert.True(f.isSet())
	assert.False(f.Empty())

	// Test empty

	v = f.getScanTarget()
	*v.(**int) = nil
	assert.False(f.Empty())

	f.updateData(v)

	assert.Equal(0, value)
	assert.True(f.Empty())
	assert.False(f.isSet())
}
