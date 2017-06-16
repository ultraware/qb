package main

var codeTemplate = `///// {{.Table}} /////
var qb{{.Table}}Table = qb.Table{Name: "{{.TableString}}"}

{{range .Fields -}}
var qb{{$.Table}}F{{.Name}} = qb.TableField{Parent: &qb{{$.Table}}Table, Name: "{{.String}}", 
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
func (t *{{.Table}}Type) Select(f ...qb.DataField) qb.SelectBuilder {
	return t.table.Select(f...)
}

// Delete creates a DELETE query
func (t *{{.Table}}Type) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts an UPDATE query
func (t *{{.Table}}Type) Update() qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts an INSERT query
func (t *{{.Table}}Type) Insert() *qb.InsertBuilder {
	return t.table.Insert(t.All())
}

// {{.Table}} returns a new {{.Table}}Type
func {{.Table}}() *{{.Table}}Type {
	table := qb{{$.Table}}Table
	data := {{.Table}}Data{}
	return &{{.Table}}Type{
		&data,
	{{- range .Fields}}
		qb{{$.Table}}F{{.Name}}.Copy(&table).New(&data.{{.Name}}),
	{{- end}}
		&table,
	}
}
`
