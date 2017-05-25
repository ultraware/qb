package main

var codeTemplate = `///// {{.Table}} /////
var qb{{.Table}}Table = qb.Table{Name: "{{.TableString}}"}

{{range .Fields -}}
var qb{{$.Table}}F{{.Name}} = qb.TableField{Parent: &qb{{$.Table}}Table, Name: "{{.String}}", Type: "{{.Type}}", 
	{{- if .ReadOnly }}ReadOnly: true,{{end -}}
	{{- if .HasDefault }}HasDefault: true,{{end -}}
	{{- if .Primary }}Primary: true{{end -}}
}
{{end}}

// {{.Table}}Data represents the data of a single row of table "{{.Table}}"
type {{.Table}}Data struct {
	{{- range .Fields}}
		{{.Name}} {{.FieldType}}
	{{- end}}
	}

// {{.Table}}Type represents the table "{{.Table}}"
type {{.Table}}Type struct {
	Data *{{.Table}}Data
{{- range .Fields}}
	{{.Name}} qb.DataField
{{- end}}
	*qb.Table
}

// All returns every field as an array
func (t *{{.Table}}Type) All() []qb.DataField {
	return []qb.DataField{
		{{- range .Fields -}}
			t.{{.Name}},
		{{- end -}}
	}
}

// {{.Table}} returns a new {{.Table}}Type
func {{.Table}}() *{{.Table}}Type {
	table := qb{{$.Table}}Table
	data := {{.Table}}Data{}
	return &{{.Table}}Type{
		&data,
	{{- range .Fields}}
		qb.NewDataField(qb{{$.Table}}F{{.Name}}.New(&table), &data.{{.Name}}),
	{{- end}}
		&table,
	}
}
`
