package pgarchitect

import (
	"database/sql"
	"strings"

	"git.ultraware.nl/Ultraware/qb/v2"
	"git.ultraware.nl/Ultraware/qb/v2/driver/pgqb"
	"git.ultraware.nl/Ultraware/qb/v2/qb-architect/internal/db"
	"git.ultraware.nl/Ultraware/qb/v2/qb-architect/internal/db/pgarchitect/pgmodel"
	"git.ultraware.nl/Ultraware/qb/v2/qb-architect/internal/util"
	"git.ultraware.nl/Ultraware/qb/v2/qbdb"
	"git.ultraware.nl/Ultraware/qb/v2/qc"
	"git.ultraware.nl/Ultraware/qb/v2/qf"

	// pgsql driver
	_ "github.com/lib/pq"
)

// Driver implements db.Driver
type Driver struct {
	DB qbdb.DB
}

// New opens a database connection and returns a Driver
func New(dsn string) db.Driver {
	d, err := sql.Open("postgres", dsn)
	util.PanicOnErr(err)

	return Driver{pgqb.New(d)}
}

func schemas() qb.Field {
	return qf.NewCalculatedField(`current_schemas(false)`)
}

func regtype(f qb.Field) qb.Field {
	return qf.NewCalculatedField(f, `::regtype`)
}

func attrelid(schema, table string) qb.Field {
	return qf.NewCalculatedField(qb.Value(schema+`.`+table), `::regclass`)
}

func castAny(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`ANY(`, f, `)`)
}

// GetTables returns all tables in the database
func (d Driver) GetTables() []string {
	it := pgmodel.Tables()

	q := it.Select(it.TableSchema, it.TableName).
		Where(qc.Eq(it.TableSchema, castAny(schemas()))).
		GroupBy(it.TableSchema, it.TableName)

	rows, err := d.DB.Query(q)
	util.PanicOnErr(err)

	var tables []string
	for rows.Next() {
		var schema, table string
		err := rows.Scan(&schema, &table)
		util.PanicOnErr(err)

		tables = append(tables, schema+`.`+table)
	}

	return tables
}

// GetFields returns all fields in a table
func (d Driver) GetFields(table string) []db.Field {
	sp := strings.Split(table, `.`)
	schema := sp[0]
	table = sp[1]

	pa := pgmodel.PgAttribute()
	c := pgmodel.Columns()

	q := pa.Select(
		pa.Attname, regtype(pa.Atttypid), pa.Attnotnull,
		qf.Case().
			When(qc.Gt(pa.Attlen, -1), pa.Attlen).
			When(qc.NotNull(c.CharacterMaximumLength), c.CharacterMaximumLength).
			Else(0),
	).
		InnerJoin(c.ColumnName, pa.Attname,
			qc.Eq(c.TableName, table),
			qc.Eq(c.TableSchema, schema),
		).
		Where(
			qc.Eq(pa.Attrelid, attrelid(schema, table)),
			qc.Gt(pa.Attnum, 0),
			qc.Eq(pa.Attisdropped, false),
		).
		GroupBy(pa.Attname, pa.Atttypid, pa.Attnotnull, pa.Attlen, c.CharacterMaximumLength)

	rows, err := d.DB.Query(q)
	if err != nil {
		panic(err)
	}

	var fields []db.Field
	for rows.Next() {
		f := db.Field{}

		var notNullable bool

		err := rows.Scan(&f.Name, &f.Type, &notNullable, &f.Size)
		util.PanicOnErr(err)

		f.Nullable = !notNullable

		fields = append(fields, f)
	}

	return fields
}
