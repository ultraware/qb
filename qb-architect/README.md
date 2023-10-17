# qb-architect

qb-architect is used to generate `db.json` for qb.

Assume a database `example` with table:

```SQL
CREATE TYPE e_type AS ENUM (
	'a',
	'b'
);

CREATE TABLE example (
	pk serial      PRIMARY KEY,
	t  varchar(50) NOT NULL,
	i  integer     NOT NULL,
	e  e_type      NOT NULL,
	b  boolean     NOT NULL,
	n  int         NULL
);
```

This database will render this `db.json`:

```json
[
  {
    "name": "public.example",
    "fields": [
      {
        "name": "b",
        "data_type": "boolean",
        "size": 1
      },
      {
        "name": "e",
        "data_type": "e_type",
        "size": 4
      },
      {
        "name": "i",
        "data_type": "integer",
        "size": 4
      },
      {
        "name": "n",
        "data_type": "integer",
        "null": true,
        "size": 4
      },
      {
        "name": "pk",
        "data_type": "integer",
        "size": 4
      },
      {
        "name": "t",
        "data_type": "character varying",
        "size": 50
      }
    ]
  }
]
```

qb-architect it the tool that will generate this output for you

## Example useage

To generate a db.json, qb-architect needs to know which database to generate the
json from and the connection string

```bash
$ qb-architect -help
Usage of qb-architect:
  -dbms string
    	Database type to use: psql, mysql, mssql
  -fexclude value
    	Regular expressions to exclude fields
  -fonly value
    	Regular expressions to whitelist fields, only tables that match at least one are returned
  -texclude value
    	Regular expressions to exclude tables
  -tonly value
    	Regular expressions to whitelist tables, only tables that match at least one are returned
```

To render the above example the next command will be used

```bash
$ qb-architect -dbms psql "host=/tmp dbname=example"
```

#### Exclude and only

If for some reason you would not want or only want a table `t` or field `f` you can
use the exclude and only options

```bash
  -fexclude value
        Regular expressions to exclude fields
  -fonly value
        Regular expressions to whitelist fields, only tables that match at least one are returned
  -texclude value
        Regular expressions to exclude tables
  -tonly value
        Regular expressions to whitelist tables, only tables that match at least one are returned
```
