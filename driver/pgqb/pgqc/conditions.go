package pgqc // import "github.com/Ultraware/qb/v3/driver/pgqb/pgqc"

import (
	"github.com/Ultraware/qb/v3"
	"github.com/Ultraware/qb/v3/qc"
)

// ILike is a postgres-specific version of qc.Like
func ILike(f1 qb.Field, s string) qb.Condition {
	f2 := qb.MakeField(s)
	return qc.NewCondition(f1, ` ILIKE `, f2)
}
