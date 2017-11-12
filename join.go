package qb

type join struct {
	Type       Join
	New        Source
	Conditions []Condition
}

// Join is the type of join
type Join string

// All possible join types
const (
	JoinInner Join = `INNER`
	JoinLeft       = `LEFT`
	JoinRight      = `RIGHT`
	JoinCross      = `CROSS`
)
