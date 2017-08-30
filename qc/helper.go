package qc

import "git.ultraware.nl/NiseVoid/qb"

func createCondition(values ...interface{}) qb.Condition {
	return func(c *qb.Context) string {
		return qb.ConcatQuery(c, values...)
	}
}
