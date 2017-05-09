package qb

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

// UpdateExcluded generates an UPDATE statement using EXCLUDED.fieldname for every field
func UpdateExcluded(f []DataField) string {
	b := sqlBuilder{&NoAlias{}, nil}
	return `UPDATE SET ` + b.SetExcluded(f)
}

func buildUpdate(t *Table, f []DataField) (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}

	p := GetPrimaryFields(f)
	if len(p) == 0 {
		panic(`Cannot update without primary`)
	}

	return b.Update(t) + b.Set(f) + b.WhereDataField(p), b.values
}
