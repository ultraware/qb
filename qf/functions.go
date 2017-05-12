package qf

import "git.ultraware.nl/NiseVoid/qb"

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

// Lower ...
func Lower(f qb.Field) qb.Field {
	return newCalculatedField(f.Source(), `string`, `lower(`, f, `)`)
}

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
