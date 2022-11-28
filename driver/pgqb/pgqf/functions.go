package pgqf // import "git.ultraware.nl/Ultraware/qb/driver/pgqb/pgqf"

import (
	"git.ultraware.nl/Ultraware/qb"
	"git.ultraware.nl/Ultraware/qb/qf"
)

// Excluded is a postgres-specific version of qf.Excluded
func Excluded(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`EXCLUDED.`, f)
}

// ArrayAgg is the postgres-specific function array_agg
func ArrayAgg(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`array_agg(`, f, `)`)
}
