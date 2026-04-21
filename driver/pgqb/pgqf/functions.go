package pgqf // import "github.com/Ultraware/qb/v3/driver/pgqb/pgqf"

import (
	"github.com/Ultraware/qb/v3"
	"github.com/Ultraware/qb/v3/qf"
)

// Excluded is a postgres-specific version of qf.Excluded
func Excluded(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`EXCLUDED.`, f)
}

// ArrayAgg is the postgres-specific function array_agg
func ArrayAgg(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`array_agg(`, f, `)`)
}
