package test

import (
	"reflect"
	"runtime"
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qf"

	"github.com/stretchr/testify/assert"
)

func TestFunctions(t *testing.T) {
	f := qb.TableField{Name: `a`}

	checkFunction(t, f, qf.Distinct, `DISTINCT a.a`)
	checkFunction(t, f, qf.Count, `count(a.a)`)
	checkFunction(t, f, qf.Sum, `sum(a.a)`)
	checkFunction(t, f, qf.Average, `avg(a.a)`)
	checkFunction(t, f, qf.Min, `min(a.a)`)
	checkFunction(t, f, qf.Max, `max(a.a)`)
}

func checkFunction(t *testing.T, field qb.Field, f func(qb.Field) qb.Field, s string) {
	assert := assert.New(t)
	fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	t.Log(`Testing`, fn)
	assert.Equal(s, f(field).QueryString(`a`), fn+` failed`)
}
