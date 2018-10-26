package db

// Driver is the minimal interface needed for db.json generation
type Driver interface {
	GetTables() []string
	GetFields(string) []Field
}

// Table represents the basic structure of db.json
type Table struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

// Field represents the structure for a db.json field
type Field struct {
	Name     string `json:"name"`
	Type     string `json:"data_type,omitempty"`
	Nullable bool   `json:"null,omitempty"`
	Size     int    `json:"size,omitempty"`
}
