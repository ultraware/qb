# qb

qb is a library that allows you to build queries without using strings. This offers some unique advantages:

- When changing your database queries that refer to old fields or tables won't compile until you update them
- You can't misspell keywords or fieldnames, this saves a lot of time and many bugs
- You can use tab completion
- You can easily port a query to a different database
- The order of commands in your query does not matter, this makes building queries in parts or adding optional statements easier

## Installation

```bash
go get git.ultraware.nl/NiseVoid/qb/...
```

## Quick start guide

### 1. Create a db.json

You can create a db.json manually or use qb-architect to generate it from your database

`qb-architect` example:

```bash
qb-architect -dbms psql host=127.0.0.1 username=qb_test dbname=qb_test > db.json
```

`db.json` example:

```json
[
	{
		"name": "TableOne",
		"alias": "one", // optional
		"fields": [
			{
				"name": "Field1",
				"data_type": "int",     // optional
				"read_only": true       // optional
			},
			{
				"name": "Field2",
				"data_type": "varchar", // optional
				"size": 50,             // optional
			},
			{ ... }
		]
	},
	{
		"name": "TableTwo",
		"fields": [
			{"name": "Field1"},
			{"name": "Field2"},
			{"name": "Field3"}
		]
	}
]
```

### 2. Run qb-generator

```bash
qb-generator db.json tables.go
```

#### Recommendations

- Don't commit qb-generator's generated code to your repo
- Use a go generate command to run qb-generator

### 3. Make a qbdb.DB

```golang
package main

var db *qbdb.DB

func main() {
	database, err := sql.Open(driver, connectionString)
	if err != nil {
		panic(err)
	}

	db = autoqb.New(database)
}
```

### 4. Write queries!

You can now write queries, you can find examples below

## Examples

### Select

```golang
one := model.One()

q := one.Select(one.Field1, one.Field2).
	Where(qc.In(Field1, 1, 2, 3))

rows, err := db.Query(q)
if err != nil {
	panic(err)
}

for rows.Next() {
	f1, f2 := 0, ""
	err := rows.Scan(&f1, &f2)
	if err != nil {
		panic(err)
	}

	fmt.Println(f1, f2)
}
```

### Insert

```golang
one := model.One()

q := one.Insert(one.Field1, one.Field2).
	Values(1, "Record 1").
	Values(2, "Record 2").
	Values(4, "Record 4")

_, err := db.Exec(q)
if err != nil {
	panic(err)
}
```

### Update

```golang
one := model.One()

q := one.Update().
	Set(one.Field2, "Record 3").
	Where(qc.Eq(one.Field1, 4))

_, err := db.Exec(q)
if err != nil {
	panic(err)
}
```

### Delete

```golang
one := model.One()

q := one.Delete(qc.Eq(one.Field1, 4))

_, err := db.Exec(q)
if err != nil {
	panic(err)
}
```

### Prepare

```golang
one := model.One()

id := 0
q := one.Select(one.Field1, one.Field2).
	Where(qc.Eq(one.Field1, &id))

stmt, err := db.Prepare()
if err != nil {
	panic(err)
}

for _, v := range []int{1,2,3,4,5} {
	id = v

	row := stmt.QueryRow()

	f1, f2 := 0, ""
	err := row.Scan(&field1, &field2)
	if err != nil {
		panic(err)
	}

	fmt.Println(f1, f2)
}
```

### Subqueries

```golang
one := model.One()

sq := one.Select(one.Field1).SubQuery()

q := sq.Select(sq.F[0])

rows, err := db.Query(q)
if err != nil {
	panic(err)
}

for rows.Next() {
	f1 := 0
	err := rows.Scan(&f1)
	if err != nil {
		panic(err)
	}

	fmt.Println(f1)
}
```

Alternatively, `.CTE()` can be used instead of `.SubQuery()` to use a CTE instead of a subquery

### Custom functions

```golang
func dbfunc(f qb.Field) qb.Field {
    return qf.NewCalculatedField("dbfunc(", f, ")")
}
```

```golang
q := one.Select(dbfunc(one.Field1))
```

### Custom conditions

```golang
func dbcond(f qb.Field) qb.Condition {
	return qc.NewCondition("dbcond(", f, ")")
}
```

```golang
q := one.Select(one.Field1).
	Where(dbcond(one.Field1))
```
