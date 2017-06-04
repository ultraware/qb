package qb

// InsertHeaderSQL ...
func InsertHeaderSQL(t *Table, f []DataField) string {
	b := sqlBuilder{&NoAlias{}, nil}
	return `INSERT INTO ` + t.QueryString(b.alias, &b.values) + ` (` + b.ListDataFields(f, false) + `) VALUES` + "\n"
}

// InsertValueSQL ...
func InsertValueSQL(f []DataField) (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}
	var s string

	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		if shouldDefault(v) {
			s += `DEFAULT`
			continue
		}
		s += Value(v.GetValue()).QueryString(b.alias, &b.values)
	}
	s = `(` + s + `)`
	return s, b.values
}

func shouldDefault(v DataField) bool {
	field, ok := v.Field.(*TableField)
	if !ok {
		panic(`Cannot use non-table field in insert`)
	}
	return (!v.isSet() || field.ReadOnly) && field.HasDefault
}
