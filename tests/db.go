package tests

import (
	"database/sql"
	"strings"

	"git.ultraware.nl/NiseVoid/qb/driver/msqb"
	"git.ultraware.nl/NiseVoid/qb/driver/myqb"
	"git.ultraware.nl/NiseVoid/qb/driver/pgqb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"

	_ "github.com/denisenkom/go-mssqldb" // database driver
	_ "github.com/go-sql-driver/mysql"   // database driver
	_ "github.com/lib/pq"                // database driver
)

func initDatabase(driverName, connectionString string, driverFunc func(*sql.DB) *qbdb.DB) *qbdb.DB {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		panic(err)
	}

	_, _ = db.Exec(`DROP TABLE one, two`)

	sql := createSQL
	if driverName != `postgres` {
		sql = strings.Replace(sql, `timestamp`, `datetime`, -1)
	}

	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	return driverFunc(db)
}

func initPostgres() *qbdb.DB {
	return initDatabase(`postgres`, getPostgresDBString(), pgqb.New)
}

func initMysql() *qbdb.DB {
	return initDatabase(`mysql`, getMysqlDBString(), myqb.New)
}

func initMssql() *qbdb.DB {
	return initDatabase(`mssql`, getMssqlDBString(), msqb.New)
}
