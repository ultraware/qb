package tests

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/autoqb"
	"git.ultraware.nl/NiseVoid/qb/internal/tests/internal/model"
	"git.ultraware.nl/NiseVoid/qb/internal/testutil"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"
)

var db *qbdb.DB
var driver string

func TestEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var db *sql.DB
	switch os.Getenv(`TYPE`) {
	case `postgres`:
		db = initPostgres()
	case `mysql`:
		db = initMysql()
	case `mssql`:
		db = initMssql()
	default:
		t.Skip(`Missing TYPE`)
	}

	startTests(t, db)
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

	testRollback(t)

	testSelect(t)
	testSelectOffset(t)
	testInQuery(t)
	testExists(t)
	testPrepare(t)
	testSubQuery(t)
	testUnionAll(t)
	if driver == `myqb.Driver` {
		testDelete(t)
	} else {
		testDeleteReturning(t)
	}
	testLeftJoin(t)
}

func testUpsert(test *testing.T) {
	o := model.One()

	q := o.Insert(o.ID, o.Name).
		Values(1, `Test 1`).
		Values(2, `Test 2`)
	q.Upsert(
		o.Update().
			Set(o.Name, qf.Concat(qf.Excluded(o.Name), `.1`)),
		o.ID,
	)
	res := db.MustExec(q)

	assert := assert.New(test)
	assert.True(res.MustRowsAffected() >= 2)
}

func testUpsertSeperate(test *testing.T) {
	o := model.One()

	iq := o.Insert(o.ID, o.Name).
		Values(1, `Test 1`).
		Values(2, `Test 2`)
	res, err := db.Exec(iq)

	assert := assert.New(test)
	if err == nil {
		assert.Eq(int64(2), res.MustRowsAffected())
		return
	}

	uq := o.Update().
		Set(o.Name, qf.Concat(o.Name, `.1`))

	res = db.MustExec(uq)

	assert.Eq(int64(2), res.MustRowsAffected())
}

func testInsert(test *testing.T) {
	t := model.Two()

	tx := db.MustBegin()

	q := t.Insert(t.OneID, t.Number, t.Comment).
		Values(1, 1, `Test comment`).
		Values(1, 2, `Test comment 2`)

	res := tx.MustExec(q)

	assert := assert.New(test)
	assert.Eq(int64(2), res.MustRowsAffected())

	tx.MustCommit()
}

func testUpdate(test *testing.T) {
	t := model.Two()

	q := t.Update().
		Set(t.Comment, qf.Concat(t.Comment, ` v2`)).
		Where(qc.Eq(t.OneID, 1))

	res := db.MustExec(q)

	assert := assert.New(test)
	assert.Eq(int64(2), res.MustRowsAffected())
}

func testUpdateReturning(test *testing.T) {
	t := model.Two()

	q := t.Update().
		Set(t.Comment, qf.Concat(t.Comment, ` v2`)).
		Where(qc.Eq(t.OneID, 1))

	r := db.MustQuery(qb.Returning(q, t.Comment, t.Number))

	assert := assert.New(test)

	var (
		comment = ``
		number  = 0
	)

	assert.True(r.Next())
	assert.NoError(r.Scan(&comment, &number))
	assert.Eq(`Test comment v2`, comment)
	assert.Eq(1, number)

	assert.True(r.Next())
	r.MustScan(&comment, &number)
	assert.Eq(`Test comment 2 v2`, comment)
	assert.Eq(2, number)

	assert.False(r.Next())
}

func testSelect(test *testing.T) {
	o, t := model.One(), model.Two()

	q := o.Select(o.ID, o.Name, qf.Year(o.CreatedAt), t.Number, t.Comment, t.ModifiedAt).
		InnerJoin(t.OneID, o.ID).
		Where(qc.Eq(o.ID, 1))
	r := db.QueryRow(q)

	var (
		id           int
		name         string
		year, number int
		comment      string
		modified     *time.Time
		assert       = assert.New(test)
	)

	assert.True(r.MustScan(&id, &name, &year, &number, &comment, &modified))

	assert.Eq(1, id)
	assert.Eq(`Test 1.1`, name)
	assert.Eq(time.Now().Year(), year)

	assert.Eq(1, number)
	assert.Eq(`Test comment v2`, comment)
	assert.Nil(modified)
}

func testSelectOffset(test *testing.T) {
	o := model.One()

	q := o.Select(o.ID, o.Name, qf.Year(o.CreatedAt)).
		OrderBy(qb.Asc(o.ID)).
		Limit(2).
		Offset(1)
	r := db.QueryRow(q)

	var (
		id     int
		name   string
		year   int
		assert = assert.New(test)
	)

	assert.True(r.MustScan(&id, &name, &year))

	assert.Eq(2, id)
	assert.Eq(`Test 2.1`, name)
	assert.Eq(time.Now().Year(), year)
}

func testInQuery(test *testing.T) {
	o, o2 := model.One(), model.One()

	sq := o2.Select(o2.ID).Where(qc.Eq(o2.ID, 1))

	q := o.Select(o.Name).
		Where(qc.InQuery(o.ID, sq))
	row := db.QueryRow(q)

	var name string
	assert := assert.New(test)

	assert.True(row.MustScan(&name))

	assert.Eq(`Test 1.1`, name)
}

func testExists(test *testing.T) {
	o, t, t2 := model.One(), model.Two(), model.Two()

	sq := t2.Select(t2.OneID).Where(qc.Eq(t2.OneID, o.ID), qc.Eq(t2.OneID, t.OneID))

	q := o.Select(qf.Count(qf.Distinct(o.Name)), qf.Count(t.OneID)).
		LeftJoin(t.OneID, o.ID).
		Where(qc.Exists(sq))
	row := db.QueryRow(q)

	var names, count int
	assert := assert.New(test)

	assert.True(row.MustScan(&names, &count))

	assert.Eq(1, names)
	assert.Eq(2, count)
}

func testPrepare(test *testing.T) {
	o := model.One()

	oneid := 0

	q := o.Select(o.ID).
		Where(qc.Eq(o.ID, &oneid))

	stmt, err := db.Prepare(q)
	if err != nil {
		panic(err)
	}

	assert := assert.New(test)
	out := 0

	assert.Eq(sql.ErrNoRows, stmt.QueryRow().Scan(&out))

	oneid = 1
	assert.NoError(stmt.QueryRow().Scan(&out))
	assert.Eq(oneid, out)

	oneid = 2
	assert.True(stmt.QueryRow().MustScan(&out))
	assert.Eq(oneid, out)

	oneid = 3
	assert.False(stmt.QueryRow().MustScan(&out))
}

func testSubQuery(test *testing.T) {
	o, t := model.One(), model.Two()

	var sq struct {
		One   qb.Field
		Count qb.Field
	}
	t.Select(t.OneID, qf.CountAll()).
		GroupBy(t.OneID).
		SubQuery(&sq.One, &sq.Count)

	q := o.Select(o.ID, sq.Count).
		InnerJoin(sq.One, o.ID)
	r := db.QueryRow(q)

	var id, count int
	assert := assert.New(test)

	assert.True(r.MustScan(&id, &count))

	assert.Eq(1, id)
	assert.Eq(2, count)

	t.Select(t.OneID).SubQuery() // No fields should pass
}

func testUnionAll(test *testing.T) {
	o := model.One()

	sq1 := o.Select(o.ID).Where(qc.Eq(o.ID, 1))
	sq2 := o.Select(o.ID).Where(qc.Eq(o.Name, `Test 1.1`))

	q := qb.UnionAll(sq1, sq2)
	r, err := db.Query(q)
	if err != nil {
		panic(err)
	}

	var (
		id     int
		assert = assert.New(test)
	)

	assert.True(r.Next())
	r.MustScan(&id)
	assert.Eq(1, id)

	assert.True(r.Next())
	r.MustScan(&id)
	assert.Eq(1, id)

	assert.False(r.Next())
}

func testDeleteReturning(test *testing.T) {
	t := model.Two()

	q := t.Delete(qc.Eq(t.OneID, 1))
	r, err := db.Query(qb.Returning(q, t.Number))
	if err != nil {
		panic(err)
	}

	var number int

	assert := assert.New(test)

	assert.True(r.Next())
	r.MustScan(&number)
	assert.Eq(1, number)

	assert.True(r.Next())
	r.MustScan(&number)
	assert.Eq(2, number)

	assert.False(r.Next())
}

func testDelete(test *testing.T) {
	t := model.Two()

	q := t.Delete(qc.Eq(t.OneID, 1))
	_, err := db.Exec(q)

	assert := assert.New(test)
	assert.NoError(err)
}

func testLeftJoin(test *testing.T) {
	o, t := model.One(), model.Two()

	q := o.Select(o.ID, t.OneID).
		LeftJoin(t.OneID, o.ID).
		Where(qc.Eq(o.ID, 1))
	r := db.QueryRow(q)

	var (
		id     int
		oneid  *int
		assert = assert.New(test)
	)

	assert.True(r.MustScan(&id, &oneid))
	assert.Eq(1, id)
	assert.Nil(oneid)
}

func testRollback(test *testing.T) {
	o := model.One()

	tx := db.MustBegin()

	q := o.Delete(qc.Eq(1, 1))
	tx.MustExec(q)

	assert := assert.New(test)
	assert.NoError(tx.Rollback())
}
