package main // import "git.ultraware.nl/Ultraware/qb/qb-architect"

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"git.ultraware.nl/Ultraware/qb/internal/filter"
	"git.ultraware.nl/Ultraware/qb/qb-architect/internal/db"
	"git.ultraware.nl/Ultraware/qb/qb-architect/internal/db/msarchitect"
	"git.ultraware.nl/Ultraware/qb/qb-architect/internal/db/myarchitect"
	"git.ultraware.nl/Ultraware/qb/qb-architect/internal/db/pgarchitect"
)

var (
	errNoTablesFoundInDatabase = errors.New(`no tables found in this database`)

	errString = `Please specify a %s, example:` + "\n\t" +
		`qb-architect -dbms psql "host=/tmp user=qb database=architect" > db.json`
)

func main() {
	dbms := flag.String(`dbms`, ``, `Database type to use: psql, mysql, mssql`)

	var tExclude, tOnly filter.Filters
	flag.Var(&tExclude, `texclude`, `Regular expressions to exclude tables`)
	flag.Var(&tOnly, `tonly`, `Regular expressions to whitelist tables, only tables that match at least one are returned`)

	var fExclude, fOnly filter.Filters
	flag.Var(&fExclude, `fexclude`, `Regular expressions to exclude fields`)
	flag.Var(&fOnly, `fonly`, `Regular expressions to whitelist fields, only tables that match at least one are returned`)

	flag.Parse()

	dsn := strings.Join(flag.Args(), ` `)
	if dsn == `` {
		println(fmt.Sprintf(errString, `connection string`))
		os.Exit(2)
	}

	var driver db.Driver
	switch strings.ToLower(*dbms) {
	case ``:
		println(fmt.Sprintf(errString, `dbms`))
		os.Exit(2)
	case `psql`, `postgres`, `postgresql`:
		driver = pgarchitect.New(dsn)
	case `mssql`, `sqlserver`:
		driver = msarchitect.New(dsn)
	case `mysql`:
		driver = myarchitect.New(dsn)
	default:
		println(`"` + *dbms + `" is not supported`)
		os.Exit(2)
	}

	filtered := filterTables(driver.GetTables(), tOnly, tExclude)

	tables := make([]db.Table, 0, len(filtered))
	for _, v := range filtered {
		tables = append(tables, db.Table{
			Name:   v,
			Fields: filterFields(driver.GetFields(v), fOnly, fExclude),
		})
	}

	err := output(tables)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func filterTables(tables []string, only, exclude filter.Filters) []string {
	var out []string

	for _, v := range tables {
		if applyFilters(v, only, exclude) {
			out = append(out, v)
		}
	}

	sort.Strings(out)

	return out
}

func filterFields(field []db.Field, only, exclude filter.Filters) []db.Field {
	var out []db.Field

	for _, v := range field {
		if applyFilters(v.Name, only, exclude) {
			out = append(out, v)
		}
	}

	sort.Sort(fields(out))

	return out
}

func applyFilters(name string, only, exclude filter.Filters) bool {
	pass := false
	for _, re := range only {
		if re.MatchString(name) {
			pass = true
			break
		}
	}

	if !pass && len(only) > 0 {
		return false
	}

	for _, re := range exclude {
		if re.MatchString(name) {
			return false
		}
	}

	return true
}

func output(tables []db.Table) error {
	if len(tables) == 0 {
		return errNoTablesFoundInDatabase
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")

	return enc.Encode(tables)
}
