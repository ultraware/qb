package tests

import (
	"database/sql"
	"strings"

	"git.ultraware.nl/NiseVoid/qb/driver/myqb"
	"git.ultraware.nl/NiseVoid/qb/driver/pgqb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"

	_ "github.com/go-sql-driver/mysql" // database driver
	_ "github.com/lib/pq"              // database driver
)

func initPostgres() *qbdb.DB {
	db, err := sql.Open(`postgres`, getPostgresDBString())
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createSQL)
	if err != nil {
		panic(err)
	}

	return pgqb.New(db)
}

func initMysql() *qbdb.DB {
	db, err := sql.Open(`mysql`, getMysqlDBString())
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(strings.Replace(createSQL, `timestamp`, `datetime`, -1))
	if err != nil {
		panic(err)
	}

	return myqb.New(db)
}
