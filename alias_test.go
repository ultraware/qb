package qb

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestAliasGeneratorNil(t *testing.T) {
	assert := assert.New(t)
	ag := AliasGenerator()

	assert.Eq(``, ag.Get(nil))
}

func TestAliasGeneratorCache(t *testing.T) {
	assert := assert.New(t)
	ag := AliasGenerator()
	tbl1, tbl2 := &Table{Name: `abc`}, &Table{Name: `abcd`}

	assert.Eq(`a`, ag.Get(tbl1))
	assert.Eq(`a2`, ag.Get(tbl2))
	assert.Eq(`a`, ag.Get(tbl1))
}

func TestAliasGeneratorDuplicateTables(t *testing.T) {
	assert := assert.New(t)
	ag := AliasGenerator()
	tbl1, tbl2 := &Table{Name: `abc`}, &Table{Name: `abc`}

	assert.Eq(`a`, ag.Get(tbl1))
	assert.Eq(`a2`, ag.Get(tbl2))
}

func TestAliasGeneratorDefinedAlias(t *testing.T) {
	assert := assert.New(t)
	ag := AliasGenerator()
	tbl1, tbl2 := &Table{Name: `abc`, Alias: `q`}, &Table{Name: `abc`, Alias: `q`}

	assert.Eq(`q`, ag.Get(tbl1))
	assert.Eq(`q2`, ag.Get(tbl2))
	assert.Eq(`q`, ag.Get(tbl1))
}

func TestNoAlias(t *testing.T) {
	assert := assert.New(t)
	ag := NoAlias()
	tbl1 := &Table{Name: `abc`}

	assert.Eq(``, ag.Get(tbl1))
}
