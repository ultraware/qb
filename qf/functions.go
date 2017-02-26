package qf

import "git.ultraware.nl/NiseVoid/qb"

// Distinct ...
func Distinct(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `DISTINCT ` + s }, Field: f}
}

// Count ...
func Count(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `count(` + s + `)` }, Field: f}
}

// Sum ...
func Sum(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `sum(` + s + `)` }, Field: f}
}

// Average ...
func Average(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `avg(` + s + `)` }, Field: f}
}

// Min ...
func Min(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `min(` + s + `)` }, Field: f}
}

// Max ...
func Max(f qb.Field) qb.Field {
	return qb.CalculatedField{Action: func(s string) string { return `max(` + s + `)` }, Field: f}
}
