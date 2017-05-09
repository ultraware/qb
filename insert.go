package qb

func buildInsert(t *Table, f []DataField) string {
	b := sqlBuilder{&NoAlias{}, nil}
	return `INSERT INTO ` + t.QueryString(b.alias, &b.values) + ` (` + b.ListDataFields(f, false) + `) VALUES` + "\n"
}

func getInsertValue(f []DataField) (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}
	s := `(`

	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		if shouldDefault(v) {
			s += `DEFAULT`
			continue
		}
		val, _ := v.Value()
		s += Value(val).QueryString(b.alias, &b.values)
	}
	s += `)`
	return s, b.values
}

func shouldDefault(v DataField) bool {
	field, ok := v.getField().(*TableField)
	if !ok {
		panic(`Cannot use non-table field in insert`)
	}
	return (!v.isSet() || field.ReadOnly) && field.HasDefault
}
