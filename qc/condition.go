package qc

import "git.ultraware.nl/NiseVoid/qb"

// Equal ...
func Equal(f1, f2 qb.Field) qb.Condition {
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: func(f ...string) string { return f[0] + ` = ` + f[1] }}
}

// Greater ...
func Greater(f1, f2 qb.Field) qb.Condition {
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: func(f ...string) string { return f[0] + ` > ` + f[1] }}
}

// GreaterEqual ...
func GreaterEqual(f1, f2 qb.Field) qb.Condition {
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: func(f ...string) string { return f[0] + ` >= ` + f[1] }}
}

// Less ...
func Less(f1, f2 qb.Field) qb.Condition {
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: func(f ...string) string { return f[0] + ` < ` + f[1] }}
}

// LessEqual ...
func LessEqual(f1, f2 qb.Field) qb.Condition {
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: func(f ...string) string { return f[0] + ` <= ` + f[1] }}
}
