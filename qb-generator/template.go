package main

var codeTemplate = `///// {{.Table}} /////
var (
	qb{{.Table}}Table = qb.Table{Name: "{{.TableString}}"{{- if .Alias }}, Alias: "{{.Alias}}"{{end -}}}

{{range .Fields -}}
	qb{{$.Table}}F{{.Name}} = qb.TableField{Parent: &qb{{$.Table}}Table, Name: "{{.String}}",
	{{- if .ReadOnly }}ReadOnly: true,{{end -}}
	{{- if .DataType.Name }}Type: qb.{{.DataType.Name}},{{end -}}
	{{- if .DataType.Size }}Size: {{.DataType.Size}},{{end -}}
	{{- if .DataType.Null }}Nullable: true,{{end -}}
}
{{end}}
)

// {{.Table}}Type represents the table "{{.Table}}"
type {{.Table}}Type struct {
{{- range .Fields}}
	{{.Name}} qb.Field
{{- end}}
	table *qb.Table
}

// GetTable returns an object with info about the table
func (t *{{.Table}}Type) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *{{.Table}}Type) Select(f ...qb.Field) *qb.SelectBuilder {
	return t.table.Select(f)
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
func (t *{{.Table}}Type) Insert(f ...qb.Field) *qb.InsertBuilder {
	return t.table.Insert(f)
}

// {{.Table}} returns a new {{.Table}}Type
func {{.Table}}() *{{.Table}}Type {
	table := qb{{$.Table}}Table
	return &{{.Table}}Type{
	{{- range .Fields}}
		qb{{$.Table}}F{{.Name}}.Copy(&table),
	{{- end}}
		&table,
	}
}
`
