package pgqf

import (
	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qf"
)

// Excluded is a postgres-specific version of qf.Excluded
func Excluded(f qb.QueryStringer) qb.Field {
	return qf.NewCalculatedField(`EXCLUDED.`, f)
}
