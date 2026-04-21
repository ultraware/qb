package myqf // import "github.com/Ultraware/qb/v3/driver/myqb/myqf"

import (
	"github.com/Ultraware/qb/v3"
	"github.com/Ultraware/qb/v3/qf"
)

// Values is a mysql-specific version of qf.Excluded
func Values(f qb.Field) qb.Field {
	return qf.NewCalculatedField(`VALUES(`, f, `)`)
}
