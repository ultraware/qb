package tests

import (
	"database/sql"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" // database driver
	_ "github.com/go-sql-driver/mysql"   // database driver
	_ "github.com/lib/pq"                // database driver
)

func initDatabase(driverName, connectionString string) *sql.DB {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		panic(err)
	}

	dropQuery := `DROP TABLE one, "two $#!"`
	sql := createSQL
	if driverName != `postgres` {
		sql = strings.Replace(sql, `timestamp`, `datetime`, -1)
	}
	if driverName == `mysql` {
		sql = strings.Replace(sql, `"`, "`", -1)
		dropQuery = strings.Replace(dropQuery, `"`, "`", -1)
	}

	_, _ = db.Exec(dropQuery)
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	return db
}

func initPostgres() *sql.DB {
	return initDatabase(`postgres`, getPostgresDBString())
}

func initMysql() *sql.DB {
	return initDatabase(`mysql`, getMysqlDBString())
}

func initMssql() *sql.DB {
	return initDatabase(`sqlserver`, getMssqlDBString())
}
