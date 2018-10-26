package msarchitect

import (
	"database/sql"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/myqb"
	"git.ultraware.nl/NiseVoid/qb/qb-architect/db"
	"git.ultraware.nl/NiseVoid/qb/qb-architect/db/msarchitect/msmodel"
	"git.ultraware.nl/NiseVoid/qb/qb-architect/util"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"

	// mssql driver
	_ "github.com/denisenkom/go-mssqldb"
)

// Driver implements db.Driver
type Driver struct {
	DB *qbdb.DB
}

// New opens a database connection and returns a Driver
func New(dsn string) db.Driver {
	d, err := sql.Open(`sqlserver`, dsn)
	util.PanicOnErr(err)

	return Driver{myqb.New(d)}
}

func database() qb.Field {
	return qf.NewCalculatedField(`DB_NAME()`)
}

// GetTables returns all tables in the database
func (d Driver) GetTables() []string {
	it := msmodel.Tables()

	q := it.Select(it.TableName).
		Where(
			qc.Eq(it.TableType, `BASE TABLE`),
			qc.Eq(it.TableCatalog, database()),
		).
		GroupBy(it.TableName)

	rows, err := d.DB.Query(q)
	util.PanicOnErr(err)

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		util.PanicOnErr(err)

		tables = append(tables, table)
	}

	return tables
}

// GetFields returns all fields in a table
func (d Driver) GetFields(table string) (fields []db.Field) {
	c := msmodel.Columns()

	q := c.Select(c.ColumnName, c.DataType, c.IsNullable).
		Where(
			qc.Eq(c.TableCatalog, database()),
			qc.Eq(c.TableName, table),
		).
		GroupBy(c.ColumnName, c.DataType, c.IsNullable)

	rows, err := d.DB.Query(q)
	util.PanicOnErr(err)

	for rows.Next() {
		f := db.Field{}

		var isNullable string

		err := rows.Scan(&f.Name, &f.Type, &isNullable)
		util.PanicOnErr(err)

		f.Nullable = isNullable == `YES`

		fields = append(fields, f)
	}

	return
}
