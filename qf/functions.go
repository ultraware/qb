package qf

import "git.ultraware.nl/NiseVoid/qb"

// Distinct ...
func Distinct(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `DISTINCT ` + s }, Field: f, Type: f.DataType()}
}

// Count ...
func Count(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `count(` + s + `)` }, Field: f, Type: `int`}
}

// Sum ...
func Sum(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `sum(` + s + `)` }, Field: f, Type: f.DataType()}
}

// Average ...
func Average(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `avg(` + s + `)` }, Field: f, Type: `float`}
}

// Min ...
func Min(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `min(` + s + `)` }, Field: f, Type: f.DataType()}
}

// Max ...
func Max(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `max(` + s + `)` }, Field: f, Type: f.DataType()}
}
