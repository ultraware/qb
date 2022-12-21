package main

import (
	"fmt"
	"os"
	"path/filepath"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func getSchema(driver, targetSchema string, files ...string) (*schema.Schema, error) {
	realm := &schema.Realm{}

	var evaluator schemahcl.Evaluator
	switch driver {
	case `sqlite`:
		evaluator = sqlite.EvalHCL
	case `mysql`:
		evaluator = mysql.EvalHCL
	default:
		evaluator = postgres.EvalHCL
	}

	parser, err := parseHCLPaths(files...)
	if err != nil {
		return nil, err
	}

	if err := evaluator.Eval(parser, realm, map[string]cty.Value{}); err != nil {
		return nil, err
	}

	if len(realm.Schemas) < 1 {
		return nil, fmt.Errorf("error no schema's found")
	}

	if targetSchema != `` {
		sc, ok := realm.Schema(targetSchema)
		if !ok {
			return nil, fmt.Errorf("targeted schema '%s' not found", targetSchema)
		}

		return sc, nil
	}

	return realm.Schemas[0], nil
}

// parseHCLPaths parses the HCL files in the given paths. If a path represents a directory,
// its direct descendants will be considered, skipping any subdirectories. If a project file
// is present in the input paths, an error is returned.
func parseHCLPaths(paths ...string) (*hclparse.Parser, error) {
	p := hclparse.NewParser()
	for _, path := range paths {
		switch stat, err := os.Stat(path); {
		case err != nil:
			return nil, err
		case stat.IsDir():
			dir, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}
			for _, f := range dir {
				// Skip nested dirs.
				if f.IsDir() {
					continue
				}
				if err := mayParse(p, filepath.Join(path, f.Name())); err != nil {
					return nil, err
				}
			}
		default:
			if err := mayParse(p, path); err != nil {
				return nil, err
			}
		}
	}
	if len(p.Files()) == 0 {
		return nil, fmt.Errorf("no schema files found in: %s", paths)
	}
	return p, nil
}

func mayParse(p *hclparse.Parser, path string) error {
	if n := filepath.Base(path); filepath.Ext(n) != ".hcl" {
		return nil
	}

	switch f, diag := p.ParseHCLFile(path); {
	case diag.HasErrors():
		return diag
	case isProjectFile(f):
		return fmt.Errorf("cannot parse project file %q as a schema file", path)
	default:
		return nil
	}
}

func isProjectFile(f *hcl.File) bool {
	for _, blk := range f.Body.(*hclsyntax.Body).Blocks {
		if blk.Type == "env" {
			return true
		}
	}
	return false
}
