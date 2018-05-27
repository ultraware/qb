package qb

// DataType represents a type in a database
type DataType uint16

// All defined DataTypes
const (
	Int = iota + 1
	String
	Bool
	Float
	Date
	Time
)
