package qbdb

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
)

func TestImplements(t *testing.T) {
	assert := assert.New(t)

	_, ok := interface{}(DB{}).(Target)
	assert.True(ok)
	_, ok = interface{}(Tx{}).(Target)
	assert.True(ok)

	_, ok = interface{}(Result{}).(Result)
	assert.True(ok)
}
