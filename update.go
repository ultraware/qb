package qb

func buildUpdate(t *Table, f []DataField) (string, []interface{}) {
	ag := NoAlias{}
	vl := ValueList{}
	values := []interface{}{}

	s := `UPDATE ` + t.QueryString(&ag, &vl) + ` SET `
	comma := false
	for _, v := range f {
		if !updatable(v) {
			continue
		}
		if comma {
			s += `, `
		}
		s += v.QueryString(&ag, &vl) + ` = ?`
		val, _ := v.Value()
		values = append(values, val)
		comma = true
	}

	c := []string{}
	for _, v := range f {
		field, ok := v.getField().(*TableField)
		if !ok || !field.Primary {
			continue
		}
		val, _ := v.Value()
		c = append(c, v.QueryString(&ag, &vl)+` = `+Value(val).QueryString(&ag, &vl))
	}
	s += getWhereSQL(c)

	if len(c) == 0 {
		panic(`Cannot update without primary`)
	}

	values = append(values, vl...)
	return s, values
}

func updatable(f DataField) bool {
	v, ok := f.getField().(*TableField)
	if !ok {
		panic(`Cannot use non-table field in update`)
	}

	if v.ReadOnly || !f.hasChanged() || v.Primary {
		return false
	}
	return true
}

func getWhereSQL(c []string) string {
	s := ``
	for k, v := range c {
		if k == 0 {
			s += ` WHERE ` + v
			continue
		}
		s += ` AND ` + v
	}

	return s
}
