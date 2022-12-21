package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteTemplate(t *testing.T) {
	sc, err := getSchema(`postgres`, ``, `./schema_test.hcl`)
	if err != nil {
		t.Errorf("error reading test schema, check the TestGetSchema test: %v", err)
		t.FailNow()
	}

	buff := &bytes.Buffer{}
	noErr(t, writeTemplate(buff, sc, `dbc`))

	temp := buff.String()
	eq(t, true, strings.Contains(temp, `package dbc`))
	eq(t, true, strings.Contains(temp, "qbUsersTable = qb.Table{Name: `users`}"))
	eq(t, true, strings.Contains(temp, "qbUsersFID = qb.TableField{Parent: &qbUsersTable, Name: `id`, Type: qb.Int}"))
	eq(t, true, strings.Contains(temp, "qbUsersFUsername = qb.TableField{Parent: &qbUsersTable, Name: `username`, Type: qb.String, Size: 100}"))
	eq(t, true, strings.Contains(temp, "qbUsersFBio = qb.TableField{Parent: &qbUsersTable, Name: `bio`, Type: qb.String, Nullable: true}"))

	if t.Failed() {
		t.Log(temp)
	}
}
