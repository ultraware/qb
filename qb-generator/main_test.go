package main

import "testing"

func expectCleanName(input, expected string) func(*testing.T) {
	return func(t *testing.T) {
		actual := cleanName(input)
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
