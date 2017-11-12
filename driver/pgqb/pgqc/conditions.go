package pgqc

import (
	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qc"
)

// ILike ...
func ILike(f1 qb.Field, s string) qb.Condition {
	f2 := qb.MakeField(s)
	return qc.NewCondition(f1, ` ILIKE `, f2)
}
