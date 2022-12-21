// Application qb-atlas
//
// qb-atlas will generate the go querybuilder just like the qb-generator
// using a atlas (declaritive database) configuration.
// Learn more about atlas on their website https://atlasgo.io.
//
// Usage:
// Use the application with the following command.
//
//	qb-atlas [options] <... directory | file.hcl>
//
// Options:
//
//	-o	 	specifies where to write the output to, defaults to stdout
//	--driver	specifies the database driver to read the atlas config with, defaults to postgres
//	--pkg		specifies the package name on the generated code, defaults to main
//	--schema	specifies which schema to use, defaults to the first
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Output string
	Driver string
	Pkg    string
	Schema string
	Files  []string
}

func NewConfig() *Config {
	c := &Config{}
	flag.StringVar(&c.Output, `o`, ``, `specifies where to write the output to, defaults to stdout`)
	flag.StringVar(&c.Driver, `driver`, `postgres`, `specifies the database driver to read the atlas config with, defaults to postgres`)
	flag.StringVar(&c.Pkg, `pkg`, `main`, `specifies the database driver to read the atlas config with, defaults to postgres`)
	flag.StringVar(&c.Schema, `schema`, ``, `specifies the targeted schema, defaults to first`)
	flag.Parse()

	c.Files = flag.Args()

	return c
}

func main() {
	cfg := NewConfig()
	if err := start(cfg); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func start(cfg *Config) error {
	sc, err := getSchema(cfg.Driver, cfg.Schema, cfg.Files...)
	if err != nil {
		return fmt.Errorf("error reading atlas configuration: %v", err)
	}

	var wr io.Writer
	if cfg.Output != `` {
		file, err := os.Create(cfg.Output)
		if err != nil {
			return fmt.Errorf("error cannot open or create file '%s': %v", cfg.Output, err)
		}
		defer file.Close() // nolint: errcheck
		wr = file
	} else {
		wr = os.Stdout
	}

	if err := writeTemplate(wr, sc, cfg.Pkg); err != nil {
		return fmt.Errorf("error while writing template: %v", err)
	}

	return nil
}
