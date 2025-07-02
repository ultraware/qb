package tests

import (
	"database/sql"
	"strings"

	"git.ultraware.nl/Ultraware/qb/v2/driver/autoqb"
)

func initDatabase(driverName, connectionString string) *sql.DB {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		panic(err)
	}

	dropQuery := `DROP TABLE IF EXISTS one, "two $#!", three`
	sql := createSQL

	switch {
	case autoqb.IsPostgres(db):
		sql = strings.ReplaceAll(sql, `timestamp`, `datetime`)
	case autoqb.IsMysql(db):
		sql = strings.ReplaceAll(sql, `"`, "`")
		dropQuery = strings.ReplaceAll(dropQuery, `"`, "`")
	}

	_, err = db.Exec(dropQuery)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	return db
}

func initPostgres() *sql.DB {
	return initDatabase(`postgres`, getPostgresDBString())
}

func initPostgresX() *sql.DB {
	return initDatabase(`pgx`, getPostgresDBString())
}

func initMysql() *sql.DB {
	return initDatabase(`mysql`, getMysqlDBString())
}

func initMssql() *sql.DB {
	return initDatabase(`sqlserver`, getMssqlDBString())
}
