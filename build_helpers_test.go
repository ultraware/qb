package qb

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestMakeFieldWithField(t *testing.T) {
	assert := assert.New(t)
	f := TableField{Name: `f1`}

	assert.Eq(f, MakeField(f))
}

func TestMakeFieldWithValue(t *testing.T) {
	assert := assert.New(t)

	assert.Eq(Value(`a`), MakeField(`a`))
}

func TestConcatQueryString(t *testing.T) {
	assert := assert.New(t)

	output := ConcatQuery(nil, `a`, `, `, `b`)

	assert.Eq(`a, b`, output)
}

func TestConcatQueryField(t *testing.T) {
	assert := assert.New(t)
	f := TableField{Name: `f`}

	output := ConcatQuery(NewContext(nil, NoAlias()), f, `, `, f)

	assert.Eq(`f, f`, output)
}

func TestConcatQuerySubquery(t *testing.T) {
	assert := assert.New(t)
	sq := &SelectBuilder{source: &Table{Name: `tbl`}, fields: []Field{TableField{Name: `f1`}}}

	NEWLINE, INDENT = ``, ``
	output := ConcatQuery(NewContext(nil, NoAlias()), sq, `, `, sq)
	NEWLINE, INDENT = "\n", "\t"

	assert.Eq("(SELECT f1FROM tbl t), (SELECT f1FROM tbl t)", output)
}

func TestJoinQuery(t *testing.T) {
	assert := assert.New(t)
	f := TableField{Name: `f1`}

	output := JoinQuery(NewContext(nil, NoAlias()), `, `, []interface{}{f, f, f})

	assert.Eq(`f1, f1, f1`, output)
}
