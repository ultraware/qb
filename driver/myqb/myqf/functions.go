package myqf // import "git.ultraware.nl/Ultraware/qb/driver/myqb/myqf"

import (
	"git.ultraware.nl/Ultraware/qb"
	"git.ultraware.nl/Ultraware/qb/qf"
)

// Values is a mysql-specific version of qf.Excluded
func Values(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`VALUES(`, f, `)`)
}
