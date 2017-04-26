package qb

func buildUpdate(t *Table, f []DataField) (string, []interface{}) {
	ag := NoAlias{}
	vl := ValueList{}
	values := []interface{}{}

	s := `UPDATE ` + t.QueryString(&ag, &vl) + ` SET `
	for k, v := range GetUpdatableFields(f) {
		if k > 0 {
			s += COMMA
		}
		s += v.QueryString(&ag, &vl) + ` = ` + VALUE
		val, _ := v.Value()
		values = append(values, val)
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

// GetUpdatableFields returns a list of all updatable fields
func GetUpdatableFields(f []DataField) []DataField {
	fields := []DataField{}
	for _, v := range f {
		if !updatable(v) {
			continue
		}
		fields = append(fields, v)
	}
	return fields
}

// GetUpsertSQL ...
func GetUpsertSQL(conflict []DataField, f []DataField) string {
	sql := ``
	for k, v := range conflict {
		if k > 0 {
			sql += COMMA
		}
		sql += v.QueryString(&NoAlias{}, nil)
	}
	return ` ON CONFLICT (` + sql + `) DO ` + buildUpdateExcluded(f)
}

func buildUpdateExcluded(f []DataField) string {
	ag := NoAlias{}
	vl := ValueList{}

	s := `UPDATE SET `
	for k, v := range GetUpdatableFields(f) {
		if k > 0 {
			s += COMMA
		}
		s += v.QueryString(&ag, &vl) + ` = EXCLUDED.` + v.QueryString(&ag, &vl)
	}

	return s
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
