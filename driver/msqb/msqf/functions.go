package msqf

import (
	"github.com/ultraware/qb"
	"github.com/ultraware/qb/qf"
)

// GetDate is a mssql-specific version of qf.Now
func GetDate() qb.Field {
	return qf.NewCalculatedField(`getdate()`)
}

// Concat is a mssql-specific version of qf.Concat
func Concat(i ...interface{}) qb.Field {
	return qf.CalculatedField(func(c *qb.Context) string {
		return qb.JoinQuery(c, ` + `, i)
	})
}

// DatePart is a mssql-specific version of qf.Extract
func DatePart(f qb.Field, part string) qb.Field {
	return qf.NewCalculatedField(`DATEPART(`, part, `, `, f, `)`)
}
