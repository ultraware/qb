package myarchitect

import (
	"database/sql"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/driver/myqb"
	"git.ultraware.nl/NiseVoid/qb/qb-architect/db"
	"git.ultraware.nl/NiseVoid/qb/qb-architect/db/myarchitect/mymodel"
	"git.ultraware.nl/NiseVoid/qb/qb-architect/util"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Driver implements db.Driver
type Driver struct {
	DB *qbdb.DB
}

// New opens a database connection and returns a Driver
func New(dsn string) db.Driver {
	d, err := sql.Open(`mysql`, dsn)
	util.PanicOnErr(err)

	return Driver{myqb.New(d)}
}

func database() qb.Field {
	return qf.NewCalculatedField(`database()`)
}

// GetTables returns all tables in the database
func (d Driver) GetTables() []string {
	it := mymodel.Tables()

	q := it.Select(it.TableName).
		Where(qc.Eq(it.TableSchema, database())).
		GroupBy(it.TableSchema, it.TableName)

	rows, err := d.DB.Query(q)
	util.PanicOnErr(err)

	var tables []string
	for rows.Next() {
		var tbl string
		err := rows.Scan(&tbl)
		util.PanicOnErr(err)

		tables = append(tables, tbl)
	}

	return tables
}

// GetFields returns all fields in a table
func (d Driver) GetFields(table string) []db.Field {
	c := mymodel.Columns()

	q := c.Select(c.ColumnName, c.DataType, c.IsNullable, c.CharacterMaximumLength).
		Where(
			qc.Eq(c.TableSchema, database()),
			qc.Eq(c.TableName, table),
		).
		GroupBy(c.ColumnName, c.DataType, c.IsNullable, c.CharacterMaximumLength)

	rows, err := d.DB.Query(q)
	util.PanicOnErr(err)

	var fields []db.Field
	for rows.Next() {
		f := db.Field{}

		var isNullable string
		var characterMaximumLength sql.NullInt64

		err := rows.Scan(&f.Name, &f.Type, &isNullable, &characterMaximumLength)
		util.PanicOnErr(err)

		f.Nullable = isNullable == `YES`
		if characterMaximumLength.Valid {
			f.Size = int(characterMaximumLength.Int64)
		}

		fields = append(fields, f)
	}

	return fields
}
