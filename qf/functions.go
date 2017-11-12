package qf

import "git.ultraware.nl/NiseVoid/qb"

///// General functions /////

// Excluded ...
func Excluded(f qb.QueryStringer) qb.Field {
	return useOverride(nil, f)
}

// Cast ...
func Cast(f qb.Field, t qb.DataType) qb.Field {
	return useOverride(cast, f, t)
}
func cast(f qb.Field, t qb.DataType) qb.Field {
	return CalculatedField(func(c *qb.Context) string {
		return qb.ConcatQuery(c, `CAST(`, f, ` AS `, c.Driver.TypeName(t), `)`)
	})
}

// Distinct ...
func Distinct(f qb.Field) qb.Field {
	return useOverride(distinct, f)
}
func distinct(f qb.Field) qb.Field {
	return NewCalculatedField(`DISTINCT `, f)
}

// CountAll ...
func CountAll() qb.Field {
	return useOverride(countAll)
}
func countAll() qb.Field {
	return NewCalculatedField(`count(1)`)
}

// Count ...
func Count(f qb.Field) qb.Field {
	return useOverride(count, f)
}
func count(f qb.Field) qb.Field {
	return NewCalculatedField(`count(`, f, `)`)
}

// Sum ...
func Sum(f qb.Field) qb.Field {
	return useOverride(sum, f)
}
func sum(f qb.Field) qb.Field {
	return NewCalculatedField(`sum(`, f, `)`)
}

// Average ...
func Average(f qb.Field) qb.Field {
	return useOverride(average, f)
}
func average(f qb.Field) qb.Field {
	return NewCalculatedField(`avg(`, f, `)`)
}

// Min ...
func Min(f qb.Field) qb.Field {
	return useOverride(min, f)
}
func min(f qb.Field) qb.Field {
	return NewCalculatedField(`min(`, f, `)`)
}

// Max ...
func Max(f qb.Field) qb.Field {
	return useOverride(max, f)
}
func max(f qb.Field) qb.Field {
	return NewCalculatedField(`max(`, f, `)`)
}

// Coalesce ...
func Coalesce(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(coalesce, f1, i)
}
func coalesce(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(`coalesce(`, f1, `, `, f2, `)`)
}

///// String functions /////

// Lower ...
func Lower(f qb.Field) qb.Field {
	return useOverride(lower, f)
}
func lower(f qb.Field) qb.Field {
	return NewCalculatedField(`lower(`, f, `)`)
}

// Concat ...
func Concat(i ...interface{}) qb.Field {
	return useOverride(concat, i...)
}
func concat(i ...interface{}) qb.Field {
	return CalculatedField(func(c *qb.Context) string {
		return qb.JoinQuery(c, ` || `, i)
	})
}

// Replace ...
func Replace(f qb.Field, from, to interface{}) qb.Field {
	return useOverride(replace, f, from, to)
}
func replace(f qb.Field, from, to interface{}) qb.Field {
	f1, f2 := qb.MakeField(from), qb.MakeField(to)
	return NewCalculatedField(`replace(`, f, `, `, f1, `, `, f2, `)`)
}

// Substring ...
func Substring(f qb.Field, from, length interface{}) qb.Field {
	return useOverride(substring, f, from, length)
}
func substring(f qb.Field, from, length interface{}) qb.Field {
	f1 := qb.MakeField(from)
	if length == nil {
		return NewCalculatedField(`substring(`, f, `, `, f1, `)`)
	}
	f2 := qb.MakeField(length)
	return NewCalculatedField(`substring(`, f, `, `, f1, `, `, f2, `)`)
}

///// Date functions /////

// Now ...
func Now() qb.Field {
	return useOverride(now)
}
func now() qb.Field {
	return NewCalculatedField(`now()`)
}

// Extract retrieves a the given part of a date
func Extract(f qb.Field, part string) qb.Field {
	return useOverride(extract, f, part)
}
func extract(f qb.Field, part string) qb.Field {
	return NewCalculatedField(`EXTRACT(`, part, ` FROM `, f, `)`)
}

// Second ...
func Second(f qb.Field) qb.Field {
	return useOverride(second, f)
}
func second(f qb.Field) qb.Field {
	return Extract(f, `second`)
}

// Minute ...
func Minute(f qb.Field) qb.Field {
	return useOverride(minute, f)
}
func minute(f qb.Field) qb.Field {
	return Extract(f, `minute`)
}

// Hour ...
func Hour(f qb.Field) qb.Field {
	return useOverride(hour, f)
}
func hour(f qb.Field) qb.Field {
	return Extract(f, `hour`)
}

// Day ...
func Day(f qb.Field) qb.Field {
	return useOverride(day, f)
}
func day(f qb.Field) qb.Field {
	return Extract(f, `day`)
}

// Week ...
func Week(f qb.Field) qb.Field {
	return useOverride(week, f)
}
func week(f qb.Field) qb.Field {
	return Extract(f, `week`)
}

// Month ...
func Month(f qb.Field) qb.Field {
	return useOverride(month, f)
}
func month(f qb.Field) qb.Field {
	return Extract(f, `month`)
}

// Year ...
func Year(f qb.Field) qb.Field {
	return useOverride(year, f)
}
func year(f qb.Field) qb.Field {
	return Extract(f, `year`)
}

///// Mathmatical functions /////

// Abs ...
func Abs(f qb.Field) qb.Field {
	return useOverride(abs, f)
}
func abs(f qb.Field) qb.Field {
	return NewCalculatedField(`abs(`, f, `)`)
}

// Ceil ...
func Ceil(f qb.Field) qb.Field {
	return useOverride(ceil, f)
}
func ceil(f qb.Field) qb.Field {
	return NewCalculatedField(`ceil(`, f, `)`)
}

// Floor ...
func Floor(f qb.Field) qb.Field {
	return useOverride(floor, f)
}
func floor(f qb.Field) qb.Field {
	return NewCalculatedField(`floor(`, f, `)`)
}

// Round ...
func Round(f1 qb.Field, precision int) qb.Field {
	return useOverride(round, f1, precision)
}
func round(f1 qb.Field, precision int) qb.Field {
	f2 := qb.MakeField(precision)
	return NewCalculatedField(`round(`, f1, `, `, f2, `)`)
}

///// Mathmatical expressions /////

// Add ...
func Add(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(add, f1, i)
}
func add(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` + `, f2)
}

// Sub ...
func Sub(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(sub, f1, i)
}
func sub(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` - `, f2)
}

// Mult ...
func Mult(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(mult, f1, i)
}
func mult(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` * `, f2)
}

// Div ...
func Div(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(div, f1, i)
}
func div(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` / `, f2)
}

// Mod ...
func Mod(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(mod, f1, i)
}
func mod(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` % `, f2)
}

// Pow ...
func Pow(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(pow, f1, i)
}
func pow(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` ^ `, f2)
}
