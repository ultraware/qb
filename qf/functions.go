package qf

import "git.ultraware.nl/NiseVoid/qb"

///// General functions /////

// Excluded uses a value from the INSERT, only usable in an upsert query
func Excluded(f qb.QueryStringer) qb.Field {
	return useOverride(nil, f)
}

// Cast casts a value to a different type
func Cast(f qb.Field, t qb.DataType) qb.Field {
	return useOverride(cast, f, t)
}
func cast(f qb.Field, t qb.DataType) qb.Field {
	return CalculatedField(func(c *qb.Context) string {
		return qb.ConcatQuery(c, `CAST(`, f, ` AS `, c.Driver.TypeName(t), `)`)
	})
}

// Distinct removes all duplicate values for this field
func Distinct(f qb.Field) qb.Field {
	return useOverride(distinct, f)
}
func distinct(f qb.Field) qb.Field {
	return NewCalculatedField(`DISTINCT `, f)
}

// CountAll counts the number of rows
func CountAll() qb.Field {
	return useOverride(countAll)
}
func countAll() qb.Field {
	return NewCalculatedField(`count(1)`)
}

// Count counts the number of non-NULL values for this field
func Count(f qb.Field) qb.Field {
	return useOverride(count, f)
}
func count(f qb.Field) qb.Field {
	return NewCalculatedField(`count(`, f, `)`)
}

// Sum calculates the sum of all values in this field
func Sum(f qb.Field) qb.Field {
	return useOverride(sum, f)
}
func sum(f qb.Field) qb.Field {
	return NewCalculatedField(`sum(`, f, `)`)
}

// Average calculates the average of all values in this field
func Average(f qb.Field) qb.Field {
	return useOverride(average, f)
}
func average(f qb.Field) qb.Field {
	return NewCalculatedField(`avg(`, f, `)`)
}

// Min calculates the minimum value in this field
func Min(f qb.Field) qb.Field {
	return useOverride(min, f)
}
func min(f qb.Field) qb.Field {
	return NewCalculatedField(`min(`, f, `)`)
}

// Max calculates the maximum value in this field
func Max(f qb.Field) qb.Field {
	return useOverride(max, f)
}
func max(f qb.Field) qb.Field {
	return NewCalculatedField(`max(`, f, `)`)
}

// Coalesce returns the second argument if the first one is NULL
func Coalesce(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(coalesce, f1, i)
}
func coalesce(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(`coalesce(`, f1, `, `, f2, `)`)
}

///// String functions /////

// Lower returns the value as a lowercase string
func Lower(f qb.Field) qb.Field {
	return useOverride(lower, f)
}
func lower(f qb.Field) qb.Field {
	return NewCalculatedField(`lower(`, f, `)`)
}

// Concat combines fields and strings into a single string
func Concat(i ...interface{}) qb.Field {
	return useOverride(concat, i...)
}
func concat(i ...interface{}) qb.Field {
	return CalculatedField(func(c *qb.Context) string {
		return qb.JoinQuery(c, ` || `, i)
	})
}

// Replace replaces values in a string
func Replace(f qb.Field, from, to interface{}) qb.Field {
	return useOverride(replace, f, from, to)
}
func replace(f qb.Field, from, to interface{}) qb.Field {
	f1, f2 := qb.MakeField(from), qb.MakeField(to)
	return NewCalculatedField(`replace(`, f, `, `, f1, `, `, f2, `)`)
}

// Substring retrieves a part of a string
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

// Now retrieves the current time
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

// Second retrieves the second from a date
func Second(f qb.Field) qb.Field {
	return useOverride(second, f)
}
func second(f qb.Field) qb.Field {
	return Extract(f, `second`)
}

// Minute retrieves the minute from a date
func Minute(f qb.Field) qb.Field {
	return useOverride(minute, f)
}
func minute(f qb.Field) qb.Field {
	return Extract(f, `minute`)
}

// Hour retrieves the hour from a date
func Hour(f qb.Field) qb.Field {
	return useOverride(hour, f)
}
func hour(f qb.Field) qb.Field {
	return Extract(f, `hour`)
}

// Day retrieves the day from a date
func Day(f qb.Field) qb.Field {
	return useOverride(day, f)
}
func day(f qb.Field) qb.Field {
	return Extract(f, `day`)
}

// Week retrieves the week from a date
func Week(f qb.Field) qb.Field {
	return useOverride(week, f)
}
func week(f qb.Field) qb.Field {
	return Extract(f, `week`)
}

// Month retrieves the month from a date
func Month(f qb.Field) qb.Field {
	return useOverride(month, f)
}
func month(f qb.Field) qb.Field {
	return Extract(f, `month`)
}

// Year retrieves the year from a date
func Year(f qb.Field) qb.Field {
	return useOverride(year, f)
}
func year(f qb.Field) qb.Field {
	return Extract(f, `year`)
}

///// Mathmatical functions /////

// Abs returns the absolute value, turning all negatieve numbers into positive numbers
func Abs(f qb.Field) qb.Field {
	return useOverride(abs, f)
}
func abs(f qb.Field) qb.Field {
	return NewCalculatedField(`abs(`, f, `)`)
}

// Ceil rounds a value up
func Ceil(f qb.Field) qb.Field {
	return useOverride(ceil, f)
}
func ceil(f qb.Field) qb.Field {
	return NewCalculatedField(`ceil(`, f, `)`)
}

// Floor rounds a value down
func Floor(f qb.Field) qb.Field {
	return useOverride(floor, f)
}
func floor(f qb.Field) qb.Field {
	return NewCalculatedField(`floor(`, f, `)`)
}

// Round rounds a value to the specified precision
func Round(f1 qb.Field, precision int) qb.Field {
	return useOverride(round, f1, precision)
}
func round(f1 qb.Field, precision int) qb.Field {
	f2 := qb.MakeField(precision)
	return NewCalculatedField(`round(`, f1, `, `, f2, `)`)
}

///// Mathmatical expressions /////

// Add adds the values (+)
func Add(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(add, f1, i)
}
func add(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` + `, f2)
}

// Sub subtracts the values (-)
func Sub(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(sub, f1, i)
}
func sub(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` - `, f2)
}

// Mult multiplies the values (*)
func Mult(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(mult, f1, i)
}
func mult(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` * `, f2)
}

// Div divides the values (/)
func Div(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(div, f1, i)
}
func div(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` / `, f2)
}

// Mod gets the remainder of the division (%)
func Mod(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(mod, f1, i)
}
func mod(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` % `, f2)
}

// Pow calculates the power of a number (^)
func Pow(f1 qb.Field, i interface{}) qb.Field {
	return useOverride(pow, f1, i)
}
func pow(f1 qb.Field, i interface{}) qb.Field {
	f2 := qb.MakeField(i)
	return NewCalculatedField(f1, ` ^ `, f2)
}
