package autoqb // import "git.ultraware.nl/Ultraware/qb/driver/autoqb"

import (
	"database/sql"
	"reflect"
	"strings"

	"git.ultraware.nl/Ultraware/qb/driver/msqb"
	"git.ultraware.nl/Ultraware/qb/driver/myqb"
	"git.ultraware.nl/Ultraware/qb/driver/pgqb"
	"git.ultraware.nl/Ultraware/qb/qbdb"
)

// New automatically selects a qb driver
func New(db *sql.DB) qbdb.DB {
	switch getPkgName(db) {
	case `github.com/lib/pq`, `github.com/jackc/pgx/stdlib`:
		return pgqb.New(db)
	case `github.com/go-sql-driver/mysql`, `github.com/ziutek/mymysql/godrv`:
		return myqb.New(db)
	case `github.com/denisenkom/go-mssqldb`:
		return msqb.New(db)
	}
	panic(`Unknown database driver`)
}

func getPkgName(db *sql.DB) string {
	t := reflect.TypeOf(db.Driver())
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	path := t.PkgPath()
	parts := strings.Split(path, `/vendor/`)
	return parts[len(parts)-1]
}
