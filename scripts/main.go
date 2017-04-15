package main

import (
	"html/template"
	"os"
	"os/exec"
	"path"
)

// FieldType ...
type FieldType struct {
	TypeName string
	Type     string
	Base     string
	Nullable bool
}

func main() {
	os.Chdir(`scripts`)

	t, err := template.ParseFiles(`fields.tmpl`)
	if err != nil {
		panic(err)
	}

	types := []FieldType{
		FieldType{
			TypeName: `String`,
			Type:     `string`,
		},
		FieldType{
			TypeName: `Bool`,
			Type:     `bool`,
		},
		FieldType{
			TypeName: `Int`,
			Type:     `int`,
			Base:     `int64`,
		},
		FieldType{
			TypeName: `Int64`,
			Type:     `int64`,
		},
		FieldType{
			TypeName: `Int32`,
			Type:     `int32`,
			Base:     `int64`,
		},
		FieldType{
			TypeName: `Float64`,
			Type:     `float64`,
			Base:     `Float64`,
		},
		FieldType{
			TypeName: `Float32`,
			Type:     `float32`,
			Base:     `Float64`,
		},
		FieldType{
			TypeName: `Bytes`,
			Type:     `[]byte`,
		},
		FieldType{
			TypeName: `Time`,
			Type:     `time.Time`,
		},
	}

	for _, v := range types {
		types = append(types, FieldType{v.TypeName, v.Type, v.Base, true})
	}

	d, _ := os.Getwd()

	f, err := os.OpenFile(path.Join(d, `..`, `qb-fields.go`), os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}

	f.WriteString("package qb\n\n")

	for _, v := range types {
		t.Execute(f, v)
	}

	f.Close()

	exec.Command(`goimports`, `-w`, f.Name()).Run()
}
