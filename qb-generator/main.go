package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// InputTable ...
type InputTable struct {
	String string       `json:"name"`
	Fields []InputField `json:"fields"`
}

// InputField ...
type InputField struct {
	String   string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"null"`
	ReadOnly bool   `json:"read_only"`
	Default  bool   `json:"default"`
}

// Table ...
type Table struct {
	Table       string
	TableString string
	Fields      []Field
}

// Field ...
type Field struct {
	Name       string
	String     string
	Type       string
	FieldType  string
	ReadOnly   bool
	HasDefault bool
}

var fieldTypes = map[string]string{
	`time`:  `time.Time`,
	`bytes`: `[]byte`,
	`float`: `float64`,
}

var fullUpperList = []string{
	`acl`,
	`api`,
	`ascii`,
	`cpu`,
	`css`,
	`dns`,
	`eof`,
	`guid`,
	`html`,
	`http`,
	`https`,
	`id`,
	`ip`,
	`json`,
	`lhs`,
	`qps`,
	`ram`,
	`rhs`,
	`rpc`,
	`sla`,
	`smtp`,
	`sql`,
	`ssh`,
	`tcp`,
	`tls`,
	`ttl`,
	`udp`,
	`ui`,
	`uid`,
	`uuid`,
	`uri`,
	`url`,
	`utf8`,
	`vm`,
	`xml`,
	`xmpp`,
	`xsrf`,
	`xss`,
}

func getType(t string, null bool) string {
	p := ``
	if null {
		p = `*`
	}
	if v, ok := fieldTypes[t]; ok {
		return p + v
	}
	return p + t
}

func newField(name string, t string, nullable bool, readOnly bool, hasDefault bool) Field {
	return Field{cleanName(name), name, t, getType(t, nullable), readOnly, hasDefault}
}

func cleanName(s string) string {
	parts := strings.Split(s, `.`)
	parts = strings.Split(parts[len(parts)-1], `_`)
	for k := range parts {
		upper := false
		for _, v := range fullUpperList {
			if v == parts[k] {
				upper = true
				break
			}
		}

		if upper || len(parts[k]) <= 1 {
			parts[k] = strings.ToUpper(parts[k])
			continue
		}

		parts[k] = strings.ToUpper(string(parts[k][0])) + parts[k][1:]
	}
	return strings.Join(parts, ``)
}

var pkg string

func init() {
	log.SetFlags(0)

	flag.StringVar(&pkg, `package`, `model`, `The package name for the output file`)
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Println(`Usage: qbgenerate [options] input.json output.go`)
		os.Exit(2)
	}
}

func main() {
	in, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(`Failed to open input file. `, err)
	}

	input := []InputTable{}

	err = json.NewDecoder(in).Decode(&input)
	if err != nil {
		log.Fatal(`Failed to parse input file. `, err)
	}

	out, err := os.OpenFile(flag.Arg(1), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(`Failed to open output file. `, err)
	}

	err = generateCode(out, input)
	if err != nil {
		log.Fatal(`Failed to generate code. `, err)
	}

	_ = out.Close()
	err = exec.Command(`goimports`, `-w`, out.Name()).Run()
	if err != nil {
		log.Fatal(`Failed to exectue goimports. `, err)
	}
}

func generateCode(out io.Writer, input []InputTable) error {
	tables := make([]Table, len(input))
	for k, v := range input {
		t := &tables[k]
		t.Table = cleanName(v.String)
		t.TableString = v.String

		for _, f := range v.Fields {
			t.Fields = append(t.Fields, newField(f.String, f.Type, f.Nullable, f.ReadOnly, f.Default))
		}
	}

	t, err := template.New(`code`).Parse(codeTemplate)
	if err != nil {
		return err
	}

	_, _ = io.WriteString(out, `package `+pkg+"\n\n")
	for _, v := range tables {
		if err := t.Execute(out, v); err != nil {
			return err
		}
	}
	return nil
}
