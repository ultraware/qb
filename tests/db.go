package tests

import (
	"database/sql"

	"git.ultraware.nl/NiseVoid/qb/driver/pgqb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"

	_ "github.com/lib/pq" // database driver
)

func initDB() *qbdb.DB {
	db, err := sql.Open(`postgres`, getDBString())
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createSQL)
	if err != nil {
		panic(err)
	}

	return pgqb.New(db)
}
