package pgqc // import "git.ultraware.nl/Ultraware/qb/driver/pgqb/pgqc"

import (
	"git.ultraware.nl/Ultraware/qb"
	"git.ultraware.nl/Ultraware/qb/qc"
)

// ILike is a postgres-specific version of qc.Like
func ILike(f1 qb.Field, s string) qb.Condition {
	f2 := qb.MakeField(s)
	return qc.NewCondition(f1, ` ILIKE `, f2)
}
