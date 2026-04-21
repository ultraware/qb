package autoqb // import "github.com/Ultraware/qb/v3/driver/autoqb"

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/Ultraware/qb/v3/driver/msqb"
	"github.com/Ultraware/qb/v3/driver/myqb"
	"github.com/Ultraware/qb/v3/driver/pgqb"
	"github.com/Ultraware/qb/v3/qbdb"
)

// New automatically selects a qb driver
func New(db *sql.DB) qbdb.DB {
	switch getPkgName(db) {
	case `github.com/lib/pq`, `github.com/jackc/pgx/stdlib`:
		return pgqb.New(db)
	case `github.com/go-sql-driver/mysql`, `github.com/ziutek/mymysql/godrv`:
		return myqb.New(db)
	case `github.com/denisenkom/go-mssqldb`, `github.com/microsoft/go-mssqldb`:
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
