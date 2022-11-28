package qb

import (
	"reflect"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var override = OverrideMap{}

func TField() Field {
	return func() Field {
		return override.Field(GetFuncFrame(), field1, nil)
	}()
}

type calculatedField func(*Context) string

func (f calculatedField) QueryString(c *Context) string {
	return f(c)
}

func field1() Field {
	return calculatedField(func(_ *Context) string {
		return `test`
	})
}

func field2() Field {
	return calculatedField(func(_ *Context) string {
		return `test2`
	})
}

func TestOverrideField(t *testing.T) {
	assert := assert.New(t)

	val := reflect.ValueOf

	assert.Cmp(val(TField()), val(field1()), cmpopts.IgnoreUnexported(reflect.Value{}))

	override.Add(TField, field2)

	assert.Cmp(val(TField()), val(field2()), cmpopts.IgnoreUnexported(reflect.Value{}))
}
