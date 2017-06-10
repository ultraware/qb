package qf

import "git.ultraware.nl/NiseVoid/qb"

///// General functions /////

// Excluded ...
func Excluded(f qb.Field) *CalculatedField {
	return &CalculatedField{
		Action: func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
			return d.ExcludedField(f.QueryString(d, ag, vl))
		},
		S:    f.Source(),
		Type: f.DataType(),
	}
}

// Distinct ...
func Distinct(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `DISTINCT `, f)
}

// CountAll ...
func CountAll() *CalculatedField {
	return newCalculatedField(nil, `int`, `count(1)`)
}

// Count ...
func Count(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), `int`, `count(`, f, `)`)
}

// Sum ...
func Sum(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `sum(`, f, `)`)
}

// Average ...
func Average(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), `float`, `avg(`, f, `)`)
}

// Min ...
func Min(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `min(`, f, `)`)
}

// Max ...
func Max(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `max(`, f, `)`)
}

// Coalesce ...
func Coalesce(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), `coalesce(`, f1, `, `, f2, `)`)
}

///// String functions /////

// Lower ...
func Lower(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), `string`, `lower(`, f, `)`)
}

// Concat ...
func Concat(i ...interface{}) *CalculatedField {
	return &CalculatedField{
		Action: func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
			s := make([]interface{}, len(i)*2-1)
			for k, v := range i {
				if k > 0 {
					s[k*2-1] = ` ` + d.ConcatOperator() + ` `
				}

				f := qb.MakeField(v)

				s[k*2] = f
			}

			return qb.ConcatQuery(d, ag, vl, s...)
		},
		S:    nil,
		Type: `string`,
	}
}

///// Date functions /////

// Now ...
func Now() *CalculatedField {
	return newCalculatedField(nil, `time`, `now()`)
}

///// Mathmatical functions /////

// Abs ...
func Abs(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `abs(`, f, `)`)
}

// Ceil ...
func Ceil(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `ceil(`, f, `)`)
}

// Floor ...
func Floor(f qb.Field) *CalculatedField {
	return newCalculatedField(f.Source(), f.DataType(), `floor(`, f, `)`)
}

// Round ...
func Round(f1 qb.Field, precision int) *CalculatedField {
	f2 := qb.MakeField(precision)
	return newCalculatedField(f1.Source(), f1.DataType(), `round(`, f1, `, `, f2, `)`)
}

///// Mathmatical expressions /////

// Add ...
func Add(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` + `, f2)
}

// Sub ...
func Sub(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` - `, f2)
}

// Mult ...
func Mult(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` * `, f2)
}

// Div ...
func Div(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` / `, f2)
}

// Mod ...
func Mod(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` % `, f2)
}

// Pow ...
func Pow(f1 qb.Field, i interface{}) *CalculatedField {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` ^ `, f2)
}
