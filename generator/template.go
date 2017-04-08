package main

var codeTemplate = `// {{.Table}}
var qb{{.Table}}Table = qb.Table{Name: "{{.TableString}}"}

{{range .Fields -}}
    var qb{{$.Table}}_{{.Name}} = qb.TableField{Parent: &qb{{$.Table}}Table, Name: "{{.String}}", Type: "{{.Type}}", ReadOnly: {{if .ReadOnly}}true{{else}}false{{end}}}
{{end}}

type {{.Table}}Type struct {
{{- range .Fields}}
    {{.Name}} *qb.{{.FieldType}}
{{- end}}
    *qb.Table
}

func (t *{{.Table}}Type) All() []qb.DataField {
	return []qb.DataField{
		{{- range .Fields -}}
			t.{{.Name}},
		{{- end -}}
	}
}

func {{.Table}}() {{.Table}}Type {
    return {{.Table}}Type{
    {{- range .Fields}}
        qb.New{{.FieldType}}(&qb{{$.Table}}_{{.Name}}),
    {{- end}}
        &qb{{$.Table}}Table,
    }
}
`
