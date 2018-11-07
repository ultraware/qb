package myqf

import (
	"github.com/ultraware/qb"
	"github.com/ultraware/qb/qf"
)

// Values is a mysql-specific version of qf.Excluded
func Values(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`VALUES(`, f, `)`)
}
