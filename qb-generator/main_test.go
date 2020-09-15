package main

import (
	"regexp"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"git.ultraware.nl/NiseVoid/qb/internal/filter"
)

func expectCleanName(input, expected string) func(*testing.T) {
	return func(t *testing.T) {
		actual := cleanName(input, nil)
		if actual != expected {
			t.Errorf(`actual: "%s" does not match expected: "%s"`, actual, expected)
		}
	}
}

func TestCleanName(t *testing.T) {
	t.Run(`single space`, expectCleanName(`single table`, `SingleTable`))
	t.Run(`multiple space`, expectCleanName(`space multiple table`, `SpaceMultipleTable`))
	t.Run(`repeating space`, expectCleanName(`space  repeating  table`, `SpaceRepeatingTable`))

	t.Run(`single $`, expectCleanName(`single$table`, `SingleTable`))
	t.Run(`multiple $`, expectCleanName(`dollar$multiple$table`, `DollarMultipleTable`))
	t.Run(`repeating $`, expectCleanName(`dollar$$repeating$$table`, `DollarRepeatingTable`))

	t.Run(`single _`, expectCleanName(`single_table`, `SingleTable`))
	t.Run(`multiple _`, expectCleanName(`underscore_multiple_table`, `UnderscoreMultipleTable`))
	t.Run(`repeating _`, expectCleanName(`underscore__repeating__table`, `UnderscoreRepeatingTable`))

	t.Run(`single -`, expectCleanName(`single-table`, `SingleTable`))
	t.Run(`multiple -`, expectCleanName(`hyphen-multiple-table`, `HyphenMultipleTable`))
	t.Run(`repeating -`, expectCleanName(`hyphen--repeating--table`, `HyphenRepeatingTable`))

	t.Run(`replace all`, expectCleanName(`$makes_$no-sense at - all`, `MakesNoSenseAtAll`))
	t.Run(`preserve casing`, expectCleanName(`preServe$CASING`, `PreServeCASING`))
	t.Run(`special uppercase`, expectCleanName(`schema$base_url`, `SchemaBaseURL`))
}

func TestCleanNameTrim(t *testing.T) {
	tblRe := regexp.MustCompile(`_tbl$`)
	dollarRe := regexp.MustCompile(`.*\$`)

	cases := []struct {
		re       filter.Filters
		expected string
	}{
		{
			re:       nil,
			expected: `SchemaUserTbl`,
		},
		{
			re:       filter.Filters{tblRe},
			expected: `SchemaUser`,
		},
		{
			re:       filter.Filters{dollarRe},
			expected: `UserTbl`,
		},
		{
			re:       filter.Filters{dollarRe, tblRe},
			expected: `User`,
		},
		{ // Regexes should be executed in the correct order
			re:       filter.Filters{regexp.MustCompile(`bl$`), tblRe},
			expected: `SchemaUserT`,
		},
	}

	for _, v := range cases {
		c := v
		t.Run(c.re.String(), func(t *testing.T) {
			assert := assert.New(t)

			out := cleanName(`schema$user_tbl`, c.re)
			assert.Eq(c.expected, out)
		})
	}
}

func TestRemoveSchema(t *testing.T) {
	assert := assert.New(t)

	out := removeSchema(`public.tbl`)
	assert.Eq(`tbl`, out)

	out = removeSchema(`dbo.something.idk.tbl`)
	assert.Eq(`tbl`, out)
}
