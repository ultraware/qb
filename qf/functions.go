package qf

import "git.ultraware.nl/NiseVoid/qb"

///// General functions /////

// Distinct ...
func Distinct(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `DISTINCT `, f)
}

// CountAll ...
func CountAll() qb.Field {
	return newCalculatedField(nil, `int`, `count(1)`)
}

// Count ...
func Count(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), `int`, `count(`, f, `)`)
}

// Sum ...
func Sum(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `sum(`, f, `)`)
}

// Average ...
func Average(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), `float`, `avg(`, f, `)`)
}

// Min ...
func Min(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `min(`, f, `)`)
}

// Max ...
func Max(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `max(`, f, `)`)
}

// Coalesce ...
func Coalesce(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), `coalesce(`, f1, `, `, f2, `)`)
}

///// String functions /////

// Lower ...
func Lower(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), `string`, `lower(`, f, `)`)
}

///// Date functions /////

// Now ...
func Now() qb.Field {
	return newCalculatedField(nil, `time`, `now()`)
}

///// Mathmatical functions /////

// Abs ...
func Abs(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `abs(`, f, `)`)
}

// Ceil ...
func Ceil(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `ceil(`, f, `)`)
}

// Floor ...
func Floor(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), f.DataType(), `floor(`, f, `)`)
}

// Round ...
func Round(f1 qb.Field, precision int) qb.Field {
	f2 := makeField(precision)
	return newCalculatedField(f1.Source(), f1.DataType(), `round(`, f1, `, `, f2, `)`)
}

///// Mathmatical expressions /////

// Add ...
func Add(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` + `, f2)
}

// Sub ...
func Sub(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` - `, f2)
}

// Mult ...
func Mult(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` * `, f2)
}

// Div ...
func Div(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` / `, f2)
}

// Mod ...
func Mod(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` % `, f2)
}

// Pow ...
func Pow(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return newCalculatedField(f1.Source(), f1.DataType(), f1, ` ^ `, f2)
}
