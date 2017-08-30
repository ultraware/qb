package qf

import "git.ultraware.nl/NiseVoid/qb"

///// General functions /////

// Excluded ...
func Excluded(f qb.QueryStringer) qb.Field {
	return CalculatedField(func(c *qb.Context) string {
		return c.Driver.ExcludedField(f.QueryString(c))
	})
}

// Cast ...
func Cast(f qb.Field, t qb.DataType) qb.Field {
	return CalculatedField(func(c *qb.Context) string {
		return qb.ConcatQuery(c, `CAST(`, f, ` AS `, c.Driver.TypeName(t), `)`)
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
	return CalculatedField(func(c *qb.Context) string {
		s := make([]interface{}, len(i)*2-1)
		for k, v := range i {
			if k > 0 {
				s[k*2-1] = ` ` + c.Driver.ConcatOperator() + ` `
			}

			f := qb.MakeField(v)

			s[k*2] = f
		}

		return qb.ConcatQuery(c, s...)
	})
}

// Replace ...
func Replace(f qb.Field, from, to interface{}) qb.Field {
	f1, f2 := qb.MakeField(from), qb.MakeField(to)
	return newCalculatedField(`replace(`, f, `, `, f1, `, `, f2, `)`)
}

// Substring ...
func Substring(f qb.Field, from, length interface{}) qb.Field {
	f1 := qb.MakeField(from)
	if length == nil {
		return newCalculatedField(`substring(`, f, `, `, f1, `)`)
	}
	f2 := qb.MakeField(length)
	return newCalculatedField(`substring(`, f, `, `, f1, `, `, f2, `)`)
}

///// Date functions /////

// Now ...
func Now() qb.Field {
	return newCalculatedField(`CURRENT_TIMESTAMP`)
}

func newExtractField(f qb.Field, part string) CalculatedField {
	return func(c *qb.Context) string {
		return c.Driver.DateExtract(f.QueryString(c), part)
	}
}

// Second ...
func Second(f qb.Field) qb.Field {
	return newExtractField(f, `second`)
}

// Minute ...
func Minute(f qb.Field) qb.Field {
	return newExtractField(f, `minute`)
}

// Hour ...
func Hour(f qb.Field) qb.Field {
	return newExtractField(f, `hour`)
}

// Day ...
func Day(f qb.Field) qb.Field {
	return newExtractField(f, `day`)
}

// Week ...
func Week(f qb.Field) qb.Field {
	return newExtractField(f, `week`)
}

// Month ...
func Month(f qb.Field) qb.Field {
	return newExtractField(f, `month`)
}

// Year ...
func Year(f qb.Field) qb.Field {
	return newExtractField(f, `year`)
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
