package pgqc

import (
	"github.com/ultraware/qb"
	"github.com/ultraware/qb/qc"
)

// ILike is a postgres-specific version of qf.Like
func ILike(f1 qb.Field, s string) qb.Condition {
	f2 := qb.MakeField(s)
	return qc.NewCondition(f1, ` ILIKE `, f2)
}
