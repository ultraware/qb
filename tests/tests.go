package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/tests/internal/model"
)

var db *qbdb.DB

// StartTests starts all the end-to-end tests
func StartTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db = initDB()

	if testing.Verbose() {
		db.Debug = true
	}

	testInsert(t)
	testUpdate(t)
	testSelect(t)
	// testLeftJoin(t)
}

func testInsert(test *testing.T) {
	// Insert into One
	o := model.One()

	o.Data.ID = 1
	o.Data.Name = `Test 1`
	o.Data.CreatedAt = time.Now()

	err := db.Insert(o)
	if err != nil {
		panic(err)
	}

	// Insert into Two
	t := model.Two()

	t.Data.OneID = o.Data.ID
	t.Data.Number = 1
	t.Data.Comment = `Test comment`

	err = db.Insert(t)
	if err != nil {
		panic(err)
	}
}

func testUpdate(test *testing.T) {
	// Insert
	o := model.One()

	o.Data.ID = 1
	o.Data.Name = `Test 1.1`

	err := db.Update(o)
	if err != nil {
		panic(err)
	}

	// Insert
	t := model.Two()

	t.Data.OneID = o.Data.ID
	t.Data.Number = 1
	t.Data.Comment = `Test comment v2`

	err = db.Update(t)
	if err != nil {
		panic(err)
	}
}

func testSelect(test *testing.T) {
	o := model.One()
	t := model.Two()

	q := o.Select(o.ID, o.Name, t.Number, t.Comment, t.ModifiedAt).
		InnerJoin(t.OneID, o.ID).
		Where(qc.Eq(o.ID, 1))
	err := db.QueryRow(q)
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)

	assert.Equal(1, o.Data.ID)
	assert.False(o.ID.Empty())

	assert.Equal(`Test 1.1`, o.Data.Name)
	assert.Equal(1, t.Data.Number)
	assert.Equal(`Test comment v2`, t.Data.Comment)

	assert.Nil(t.Data.ModifiedAt)
	assert.True(t.ModifiedAt.Empty())
	assert.True(t.OneID.Empty())
}
