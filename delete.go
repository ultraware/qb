package qb

import (
	"reflect"
	"strings"
)

// DeleteBuilder builds a DELETE query
type DeleteBuilder struct {
	table *Table
	c     []Condition
}

// SQL returns a query string and a list of values
//

func (q DeleteBuilder) SQL(b SQLBuilder) (string, []interface{}) {
	drvPath := reflect.TypeOf(b.Context.Driver).PkgPath()
	if strings.HasSuffix(drvPath, `myqb`) || strings.HasSuffix(drvPath, `msqb`) {
		b.w.WriteLine(`DELETE ` + q.table.aliasString())
		b.w.WriteLine(`FROM ` + b.SourceToSQL(q.table))
	} else {
		b.Delete(q.table)
	}
	b.Where(q.c...)

	return b.w.String(), *b.Context.Values
}
