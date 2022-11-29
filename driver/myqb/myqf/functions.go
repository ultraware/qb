package myqf // import "git.ultraware.nl/Ultraware/qb/v2/driver/myqb/myqf"

import (
	"git.ultraware.nl/Ultraware/qb/v2"
	"git.ultraware.nl/Ultraware/qb/v2/qf"
)

// Values is a mysql-specific version of qf.Excluded
func Values(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`VALUES(`, f, `)`)
}
