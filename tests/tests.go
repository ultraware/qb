package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"
	"git.ultraware.nl/NiseVoid/qb/tests/internal/model"

	"github.com/fatih/color"
)

var db *qbdb.DB

// StartTests starts all the end-to-end tests
func StartTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db = initPostgres()
	runTests(t)
	db = initMysql()
	runTests(t)
}

func runTests(t *testing.T) {
	driver := reflect.TypeOf(db.Driver).String()
	if testing.Verbose() {
		db.Debug = true
		fmt.Println()
		fmt.Println(color.MagentaString(`----- Testing with: %s -----`, driver))
		fmt.Println()
	}

	testUpsert(t)
	testInsert(t)
	testUpsert(t)

	if driver == `pgqb.Driver` {
		testUpdateReturning(t)
	} else {
		testUpdate(t)
	}

	testSelect(t)
	testDelete(t)
	testLeftJoin(t)

	if testing.Verbose() {
		fmt.Println(color.MagentaString(`----- Finished testing: %s -----`, driver))
		fmt.Println()
	}
}

func testUpsert(test *testing.T) {
	o := model.One()

	o.Data.ID = 1
	o.Data.Name = `Test 1`

	q := o.Insert()
	q.Add()
	q.Upsert(
		o.Update().
			Set(o.Name, qf.Concat(qf.Excluded(o.Name), `.1`)),
		o.ID,
	)
	err := db.Exec(q)
	if err != nil {
		panic(err)
	}
}

func testInsert(test *testing.T) {
	t := model.Two()

	q := t.Insert()

	t.Data.OneID = 1
	t.Data.Number = 1
	t.Data.Comment = `Test comment`
	q.Add()

	t.Data.Number = 2
	t.Data.Comment += ` 2`
	q.Add()

	err := db.Exec(q)
	if err != nil {
		panic(err)
	}
}

func testUpdate(test *testing.T) {
	t := model.Two()

	q := t.Update().
		Set(t.Comment, qf.Concat(t.Comment, ` v2`)).
		Where(qc.Eq(t.OneID, 1))

	err := db.Exec(q)
	if err != nil {
		panic(err)
	}
}

func testUpdateReturning(test *testing.T) {
	t := model.Two()

	q := t.Update().
		Set(t.Comment, qf.Concat(t.Comment, ` v2`)).
		Where(qc.Eq(t.OneID, 1))

	cur, err := db.Query(qb.Returning(q, t.Comment, t.Number))
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)
	assert.True(cur.Next())
	assert.Equal(`Test comment v2`, t.Data.Comment)
	assert.Equal(1, t.Data.Number)
	assert.True(cur.Next())
	assert.Equal(`Test comment 2 v2`, t.Data.Comment)
	assert.Equal(2, t.Data.Number)
	assert.NoError(cur.Error())
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

func testDelete(test *testing.T) {
	t := model.Two()

	q := t.Delete(qc.Eq(t.OneID, 1))
	err := db.Exec(q)
	if err != nil {
		panic(err)
	}
}

func testLeftJoin(test *testing.T) {
	o := model.One()
	t := model.Two()

	q := o.Select(o.ID, t.OneID).
		LeftJoin(t.OneID, o.ID).
		Where(qc.Eq(o.ID, 1))
	err := db.QueryRow(q)
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)

	assert.Equal(1, o.Data.ID)
	assert.False(o.ID.Empty())

	assert.Equal(0, t.Data.OneID)
	assert.True(t.OneID.Empty())
}
