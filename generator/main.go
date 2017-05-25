package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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
	Primary  bool   `json:"primary"`
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
	Primary    bool
}

var fieldTypes = map[string]string{
	`time`:  `time.Time`,
	`bytes`: `[]byte`,
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

func newField(name string, t string, nullable bool, readOnly bool, hasDefault bool, primary bool) Field {
	return Field{cleanName(name), name, t, getType(t, nullable), readOnly, hasDefault, primary}
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

		if upper || len(parts[k]) == 0 {
			parts[k] = strings.ToUpper(parts[k])
			continue
		}

		parts[k] = strings.ToUpper(string(parts[k][0])) + parts[k][1:]
	}
	return strings.Join(parts, ``)
}

func printError(e error) {
	if e == nil {
		return
	}
	fmt.Println(e)
}

func main() {
	p := flag.String(`package`, `model`, `The package name for the output file`)
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println(`Usage: qbgenerate [options] input.json output.go`)
		os.Exit(2)
	}

	in, err := os.Open(args[0])
	if err != nil {
		fmt.Println(`Failed to open input file`)
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	input := []InputTable{}

	err = json.NewDecoder(in).Decode(&input)
	if err != nil {
		fmt.Println(`Failed to parse input file`)
		fmt.Printf("%v\n", err)
		_ = in.Close()
		os.Exit(2)
	}
	printError(in.Close())

	out, err := os.OpenFile(args[1], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(`Failed to open output file`)
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	tables := []Table{}

	for _, v := range input {
		t := Table{
			Table:       cleanName(v.String),
			TableString: v.String,
		}

		for _, f := range v.Fields {
			t.Fields = append(t.Fields, newField(f.String, f.Type, f.Nullable, f.ReadOnly, f.Default, f.Primary))
		}

		tables = append(tables, t)
	}

	fmt.Fprint(out, `package `, *p, "\n\n")

	t, err := template.New(`code`).Parse(codeTemplate)
	if err != nil {
		_ = out.Close()
		fmt.Println(`Failed to parse template`)
		fmt.Printf("%v\n", err)
		os.Exit(3)
	}

	for _, v := range tables {
		err = t.Execute(out, v)
		if err != nil {
			_ = out.Close()
			fmt.Println(`Failed to execute template`)
			fmt.Printf("%v\n", err)
			os.Exit(3)
		}
	}

	printError(out.Close())

	printError(exec.Command(`goimports`, `-w`, out.Name()).Run())
}
