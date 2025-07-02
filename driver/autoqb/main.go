package autoqb // import "git.ultraware.nl/Ultraware/qb/v2/driver/autoqb"

import (
	"database/sql"

	"git.ultraware.nl/Ultraware/qb/v2/driver/msqb"
	"git.ultraware.nl/Ultraware/qb/v2/driver/myqb"
	"git.ultraware.nl/Ultraware/qb/v2/driver/pgqb"
	"git.ultraware.nl/Ultraware/qb/v2/qbdb"
	denisenkom "github.com/denisenkom/go-mssqldb" // database driver for Microsoft MSSQL
	"github.com/go-sql-driver/mysql"              // database driver for MySQL
	"github.com/jackc/pgx/stdlib"                 // database driver for PostgreSQL
	"github.com/lib/pq"                           // database driver for PostgreSQL
	mssqldb "github.com/microsoft/go-mssqldb"     // database driver for Microsoft MSSQL
	"github.com/ziutek/mymysql/godrv"             // database driver for MySQL
)

// New automatically selects a qb driver
func New(db *sql.DB) qbdb.DB {
	switch {
	case IsPostgres(db):
		return pgqb.New(db)
	case IsMysql(db):
		return myqb.New(db)
	case IsMssql(db):
		return msqb.New(db)
	}
	panic(`Unknown database driver`)
}

func IsPostgres(db *sql.DB) bool {
	driver := db.Driver()

	if _, ok := driver.(*stdlib.Driver); ok {
		return true

	}

	_, ok := driver.(*pq.Driver)
	return ok
}

func IsMysql(db *sql.DB) bool {
	driver := db.Driver()

	if _, ok := driver.(*mysql.MySQLDriver); ok {
		return true
	}

	_, ok := driver.(*godrv.Driver)
	return ok
}

func IsMssql(db *sql.DB) bool {
	driver := db.Driver()

	if _, ok := driver.(*mssqldb.Driver); ok {
		return true
	}

	_, ok := driver.(*denisenkom.Driver)
	return ok
}
