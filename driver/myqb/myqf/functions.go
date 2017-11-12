package myqf

import (
	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qf"
)

// Values is a mysql-specific version of qf.Excluded
func Values(f qb.QueryStringer) qb.Field {
	return qf.NewCalculatedField(`VALUES(`, f, `)`)
}
