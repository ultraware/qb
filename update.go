package qb

// UpdateRecordSQL ...
func UpdateRecordSQL(t *Table, f []DataField) (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}

	p := GetPrimaryFields(f)
	if len(p) == 0 {
		panic(`Cannot update without primary`)
	}

	f = GetUpdatableFields(f)
	s := make([]set, len(f))
	for k, v := range f {
		s[k] = set{v, Value(v.Value)}
	}

	return b.Update(t) + b.Set(s) + b.WhereDataField(p), b.values
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

func updatable(f DataField) bool {
	v, ok := f.Field.(*TableField)
	if !ok {
		panic(`Cannot use non-table field in update`)
	}

	if v.ReadOnly || !f.hasChanged() || v.Primary {
		return false
	}
	return true
}

type set struct {
	Field Field
	Value Field
}

// UpdateBuilder ...
type UpdateBuilder struct {
	t   *Table
	c   []Condition
	set []set
}

// Set ...
func (q UpdateBuilder) Set(f Field, v interface{}) UpdateBuilder {
	q.set = append(q.set, set{f, makeField(v)})
	return q
}

// Where ...
func (q UpdateBuilder) Where(c Condition) UpdateBuilder {
	q.c = append(q.c, c)
	return q
}

// SQL ...
func (q UpdateBuilder) SQL() (string, []interface{}) {
	b := sqlBuilder{&NoAlias{}, nil}
	return b.Update(q.t) + b.Set(q.set) + b.Where(q.c...), b.values
}
