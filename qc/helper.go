package qc

import "git.ultraware.nl/NiseVoid/qb"

func createCondition(values ...interface{}) qb.Condition {
	return func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
		return qb.ConcatQuery(d, ag, vl, values...)
	}
}
