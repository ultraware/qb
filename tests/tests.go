package tests

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/autoqb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"
	"git.ultraware.nl/NiseVoid/qb/tests/internal/model"
	"git.ultraware.nl/NiseVoid/qb/tests/testutil"
)

var db *qbdb.DB
var driver string

// StartTests starts all the end-to-end tests
func StartTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	startTests(t, initPostgres())
	startTests(t, initMysql())
	startTests(t, initMssql())
}

func startTests(t *testing.T, d *sql.DB) {
	db = autoqb.New(d)
	driver = reflect.TypeOf(db.Driver).String()

	if testing.Verbose() {
		db.Debug = true
		fmt.Println()
		fmt.Println(testutil.Info(`Testing with:`, driver))
		fmt.Println()
	}

	runTests(t)

	if testing.Verbose() {
		fmt.Println(testutil.Info(`Finished testing:`, driver))
		fmt.Println()
	}
}

func runTests(t *testing.T) {
	if driver == `msqb.Driver` {
		testUpsertSeperate(t)
		testInsert(t)
		testUpsertSeperate(t)
	} else {
		testUpsert(t)
		testInsert(t)
		testUpsert(t)
	}

	if driver == `myqb.Driver` {
		testUpdate(t)
	} else {
		testUpdateReturning(t)
	}

	testSelect(t)
	testSubQuery(t)
	testUnionAll(t)
	testDelete(t)
	testLeftJoin(t)
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

func testUpsertSeperate(test *testing.T) {
	o := model.One()

	o.Data.ID = 1
	o.Data.Name = `Test 1`

	iq := o.Insert()
	iq.Add()
	err := db.Exec(iq)
	if err == nil {
		return
	}

	uq := o.Update().
		Set(o.Name, qf.Concat(o.Name, `.1`)).
		Where(qc.Eq(o.ID, 1))
	err = db.Exec(uq)
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
	o, t := model.One(), model.Two()

	year := 0
	q := o.Select(o.ID, o.Name, qf.Year(o.CreatedAt).New(&year), t.Number, t.Comment, t.ModifiedAt).
		InnerJoin(t.OneID, o.ID).
		Where(qc.Eq(o.ID, 1))
	err := db.QueryRow(q)
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)

	assert.Equal(1, o.Data.ID)
	assert.False(o.ID.Empty())

	assert.Equal(time.Now().Year(), year)
	assert.Equal(`Test 1.1`, o.Data.Name)
	assert.Equal(1, t.Data.Number)
	assert.Equal(`Test comment v2`, t.Data.Comment)

	assert.Nil(t.Data.ModifiedAt)
	assert.True(t.ModifiedAt.Empty())
	assert.True(t.OneID.Empty())
}

func testSubQuery(test *testing.T) {
	o, t := model.One(), model.Two()

	var count int

	sq := t.Select(t.OneID, qf.CountAll().New(&count)).
		GroupBy(t.OneID).
		SubQuery()

	q := o.Select(o.ID, sq.F[1]).
		InnerJoin(sq.F[0], o.ID)
	err := db.QueryRow(q)
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)

	assert.Equal(1, o.Data.ID)
	assert.Equal(2, count)
}

func testUnionAll(test *testing.T) {
	o := model.One()

	sq := o.Select(o.ID)

	q := qb.UnionAll(sq, sq)
	cur, err := db.Query(q)
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)

	assert.True(cur.Next())
	assert.Equal(1, o.Data.ID)
	assert.True(cur.Next())
	assert.Equal(1, o.Data.ID)

	assert.False(cur.Next())
	assert.NoError(cur.Error())
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
	o, t := model.One(), model.Two()

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
