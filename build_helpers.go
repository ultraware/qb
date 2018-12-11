package qb

// Values used when building queries
var (
	COMMA   = `, `
	NEWLINE = "\n"
	INDENT  = "\t"
	VALUE   = `?`
)

///// Field /////

// MakeField returns the value as a Field, no operation performed when the value is already a field.
// This function is not intended to be called directly
func MakeField(i interface{}) Field {
	if f, ok := i.(Field); ok {
		return f
	}
	return Value(i)
}

// ConcatQuery combines strings and Fields into string.
// This function is not intended to be called directly
func ConcatQuery(c *Context, values ...interface{}) string {
	s := ``
	for _, val := range values {
		switch v := val.(type) {
		case (Field):
			s += v.QueryString(c)
		case (SelectQuery):
			sql, _ := v.SQL(SQLBuilder{Context: c})
			s += getSubQuerySQL(sql)
		case (string):
			s += v
		}
	}
	return s
}

// JoinQuery joins fields or values into a string separated by sep.
// This function is not intended to be called directly
func JoinQuery(c *Context, sep string, values []interface{}) string {
	s := make([]interface{}, len(values)*2-1)
	for k, v := range values {
		if k > 0 {
			s[k*2-1] = sep
		}
		s[k*2] = MakeField(v)
	}

	return ConcatQuery(c, s...)
}
