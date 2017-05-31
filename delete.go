package qb

// DeleteSQL ...
func DeleteSQL(t *Table, c []Condition) (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}
	return b.Delete(t) + b.Where(c...), b.values
}

// DeleteRecordSQL ...
func DeleteRecordSQL(t *Table, f []DataField) (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}

	p := GetPrimaryFields(f)
	if len(p) == 0 {
		panic(`Cannot update without primary`)
	}

	return b.Delete(t) + b.WhereDataField(p), b.values
}
