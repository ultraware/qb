package qf

import "git.ultraware.nl/NiseVoid/qb"

///// General functions /////

// Excluded ...
func Excluded(f qb.QueryStringer) qb.Field {
	return CalculatedField(func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
		return d.ExcludedField(f.QueryString(d, ag, vl))
	})
}

// Distinct ...
func Distinct(f qb.Field) qb.Field {
	return newCalculatedField(`DISTINCT `, f)
}

// CountAll ...
func CountAll() qb.Field {
	return newCalculatedField(`count(1)`)
}

// Count ...
func Count(f qb.Field) qb.Field {
	return newCalculatedField(`count(`, f, `)`)
}

// Sum ...
func Sum(f qb.Field) qb.Field {
	return newCalculatedField(`sum(`, f, `)`)
}

// Average ...
func Average(f qb.Field) qb.Field {
	return newCalculatedField(`avg(`, f, `)`)
}

// Min ...
func Min(f qb.Field) qb.Field {
	return newCalculatedField(`min(`, f, `)`)
}

// Max ...
func Max(f qb.Field) qb.Field {
	return newCalculatedField(`max(`, f, `)`)
}

// Coalesce ...
func Coalesce(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(`coalesce(`, f1, `, `, f2, `)`)
}

///// String functions /////

// Lower ...
func Lower(f qb.Field) qb.Field {
	return newCalculatedField(`lower(`, f, `)`)
}

// Concat ...
func Concat(i ...interface{}) qb.Field {
	return CalculatedField(func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
		s := make([]interface{}, len(i)*2-1)
		for k, v := range i {
			if k > 0 {
				s[k*2-1] = ` ` + d.ConcatOperator() + ` `
			}

			f := qb.MakeField(v)

			s[k*2] = f
		}

		return qb.ConcatQuery(d, ag, vl, s...)
	})
}

///// Date functions /////

// Now ...
func Now() qb.Field {
	return newCalculatedField(`now()`)
}

///// Mathmatical functions /////

// Abs ...
func Abs(f qb.Field) qb.Field {
	return newCalculatedField(`abs(`, f, `)`)
}

// Ceil ...
func Ceil(f qb.Field) qb.Field {
	return newCalculatedField(`ceil(`, f, `)`)
}

// Floor ...
func Floor(f qb.Field) qb.Field {
	return newCalculatedField(`floor(`, f, `)`)
}

// Round ...
func Round(f1 qb.Field, precision int) qb.Field {
	f2 := qb.MakeField(precision)
	return newCalculatedField(`round(`, f1, `, `, f2, `)`)
}

///// Mathmatical expressions /////

// Add ...
func Add(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1, ` + `, f2)
}

// Sub ...
func Sub(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1, ` - `, f2)
}

// Mult ...
func Mult(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1, ` * `, f2)
}

// Div ...
func Div(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1, ` / `, f2)
}

// Mod ...
func Mod(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1, ` % `, f2)
}

// Pow ...
func Pow(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return newCalculatedField(f1, ` ^ `, f2)
}
