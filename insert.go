package qb

func buildInsert(t *Table, f []DataField) string {
	ag := NoAlias{}
	vl := ValueList{}
	s := `INSERT INTO ` + t.QueryString(&ag, &vl) + ` (`
	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		s += v.QueryString(&ag, &vl)
	}
	s += `) VALUES `
	return s
}

func getInsertValue(f []DataField) (string, []interface{}) {
	s := `(`
	values := []interface{}{}

	for k, v := range f {
		field, ok := v.getField().(*TableField)
		if !ok {
			panic(`Cannot use non-table field in insert`)
		}
		if k > 0 {
			s += COMMA
		}
		if (!v.isSet() || field.ReadOnly) && field.HasDefault {
			s += `DEFAULT`
		} else {
			s += VALUE
			val, _ := v.Value()
			values = append(values, val)
		}
	}
	s += `)`
	return s, values
}
