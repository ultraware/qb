# qb

qb is a library that allows you to build queries without using strings. This offers some unique advantages:

- When changing your database you won't be able to compile queries that refer to old fields
- You can't misspell keywords or fieldnames, this saves a lot of time
- You can use tab completion
- You can easily port a query to a different database
- The order of your query does not matter, this makes building queries in parts easier

## Installation

```bash
go get git.ultraware.nl/NiseVoid/qb/...
```

## Quick start guide

### 1. Create a db.json

You can create a db.json manually, but you can also use tools to generate this from your database

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
qb-generate db.json tables.go
```

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
	Values(1, `Record 1`).
	Values(2, `Record 2`).
	Values(4, `Record 4`)

err := db.Exec(q)
if err != nil {
	panic(err)
}
```

### Update

```golang
one := model.One()

q := one.Update().
	Set(one.Field2, `Record 3`).
	Where(qc.Eq(one.Field1, 4))

err := db.Exec(q)
if err != nil {
	panic(err)
}
```

### Delete

```golang
one := model.One()

q := one.Delete(qc.Eq(one.Field1, 4))

err := db.Exec(q)
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

for _, v := range []int{1,2,3,4,5} {
	id = v
	err := db.QueryRow(q)
	if err != nil {
		panic(err)
	}
}
```
