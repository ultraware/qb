package qbdb

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestImplements(t *testing.T) {
	assert := assert.New(t)

	_, ok := interface{}(&db{}).(DB)
	assert.True(ok)

	_, ok = interface{}(&tx{}).(Tx)
	assert.True(ok)

	_, ok = interface{}(&db{}).(Target)
	assert.True(ok)
	_, ok = interface{}(&tx{}).(Target)
	assert.True(ok)

	_, ok = interface{}(Result{}).(Result)
	assert.True(ok)
}
