package qf

import "git.ultraware.nl/NiseVoid/qb"

// Distinct ...
func Distinct(f qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `DISTINCT ` + f.QueryString(ag, vl)
		},
		S:    f.Source(),
		Type: f.DataType(),
	}
}

// CountAll ...
func CountAll() qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `count(1)`
		},
		S:    nil,
		Type: `int`,
	}
}

// Count ...
func Count(f qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `count(` + f.QueryString(ag, vl) + `)`
		},
		S:    f.Source(),
		Type: `int`,
	}
}

// Sum ...
func Sum(f qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `sum(` + f.QueryString(ag, vl) + `)`
		},
		S:    f.Source(),
		Type: f.DataType(),
	}
}

// Average ...
func Average(f qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `avg(` + f.QueryString(ag, vl) + `)`
		},
		S:    f.Source(),
		Type: `float`,
	}
}

// Min ...
func Min(f qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `min(` + f.QueryString(ag, vl) + `)`
		},
		S:    f.Source(),
		Type: f.DataType(),
	}
}

// Max ...
func Max(f qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `max(` + f.QueryString(ag, vl) + `)`
		},
		S:    f.Source(),
		Type: f.DataType(),
	}
}

// Lower ...
func Lower(f1 qb.Field) qb.Field {
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return `lower(` + f1.QueryString(ag, vl) + `)`
		},
		S:    f1.Source(),
		Type: `string`,
	}
}

// Add ...
func Add(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return f1.QueryString(ag, vl) + ` + ` + f2.QueryString(ag, vl)
		},
		S:    f1.Source(),
		Type: f1.DataType(),
	}
}

// Sub ...
func Sub(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return f1.QueryString(ag, vl) + ` - ` + f2.QueryString(ag, vl)
		},
		S:    f1.Source(),
		Type: f1.DataType(),
	}
}

// Mult ...
func Mult(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return f1.QueryString(ag, vl) + ` * ` + f2.QueryString(ag, vl)
		},
		S:    f1.Source(),
		Type: f1.DataType(),
	}
}

// Div ...
func Div(f1 qb.Field, i interface{}) qb.Field {
	f2 := makeField(i)
	return CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string {
			return f1.QueryString(ag, vl) + ` / ` + f2.QueryString(ag, vl)
		},
		S:    f1.Source(),
		Type: f1.DataType(),
	}
}
