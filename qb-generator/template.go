package main

var codeTemplate = `///// {{.Table}} /////
var qb{{.Table}}Table = qb.Table{Name: "{{.TableString}}"{{- if .Alias }}, Alias: "{{.Alias}}"{{end -}}}

{{range .Fields -}}
var qb{{$.Table}}F{{.Name}} = qb.TableField{Parent: &qb{{$.Table}}Table, Name: "{{.String}}", 
	{{- if .ReadOnly }}ReadOnly: true,{{end -}}
	{{- if .HasDefault }}HasDefault: true,{{end -}}
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
	table *qb.Table
}

// All returns every field as an array
func (t *{{.Table}}Type) All() []qb.DataField {
	return []qb.DataField{
		{{- range .Fields -}}
			t.{{.Name}},
		{{- end -}}
	}
}

// GetTable returns an object with info about the table
func (t *{{.Table}}Type) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *{{.Table}}Type) Select(f ...qb.DataField) *qb.SelectBuilder {
	return t.table.Select(f...)
}

// Delete creates a DELETE query
func (t *{{.Table}}Type) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts an UPDATE query
func (t *{{.Table}}Type) Update() *qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts an INSERT query
func (t *{{.Table}}Type) Insert() *qb.InsertBuilder {
	return t.table.Insert(t.All())
}

// {{.Table}}From returns a new {{.Table}}Type using the provided data
func {{.Table}}From(data *{{.Table}}Data) *{{.Table}}Type {
	table := qb{{$.Table}}Table
	return &{{.Table}}Type{
		data,
	{{- range .Fields}}
		qb{{$.Table}}F{{.Name}}.Copy(&table).New(&data.{{.Name}}),
	{{- end}}
		&table,
	}

}

// {{.Table}} returns a new {{.Table}}Type
func {{.Table}}() *{{.Table}}Type {
	data := {{.Table}}Data{}
	return {{.Table}}From(&data)
}
`
