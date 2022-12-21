package main

import (
	"fmt"
	"testing"
)

func TestGetSchema(t *testing.T) {
	sc, err := getSchema(`postgres`, ``, `./schema_test.hcl`)
	noErr(t, err)

	eq(t, 2, len(sc.Tables))
	for _, table := range sc.Tables {
		switch table.Name {
		case `users`:
			eq(t, 4, len(table.Columns))
			eq(t, `id`, table.Columns[0].Name)
			eq(t, `username`, table.Columns[1].Name)

		case `blogs`:
			eq(t, 6, len(table.Columns))
			eq(t, `id`, table.Columns[0].Name)
			eq(t, fmt.Sprintf("%T", table.Columns[0].Type.Type), "*postgres.SerialType")
			eq(t, `user_id`, table.Columns[1].Name)
			eq(t, 1, len(table.ForeignKeys))
			eq(t, fmt.Sprintf("%T", table.Columns[1].Type.Type), "*schema.IntegerType")

		default:
			t.Errorf("unknown table: '%s'", table.Name)
		}
	}

	// testing specified schema
	_, err = getSchema(`postgres`, `public`, `./schema_test.hcl`)
	noErr(t, err)

	// testing unknown schema
	_, err = getSchema(`postgres`, `unknown`, `./schema_test.hcl`)
	hasErr(t, err)

	// testing unknown file
	_, err = getSchema(`postgres`, ``, `./unknown_file.hcl`)
	hasErr(t, err)
}
