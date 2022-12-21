# qb-atlas

qb-atlas will generate your qb code from your [atlas configuration](https://atlasgo.io).

## Quickstart
Write your atlas configuration:
```hcl
schema "public" {
	charset = "utf-8"
}

table "users" {
	schema = schema.public

	column "id" {
		type = serial
	}

...
```

generate your code for example in the models directory while the atlas files are in the schemas directory:
```sh
~ qb-atlas -o ./models/models.gen.go -pkg models -schema public ./schemas
```


## Write queries!

You can now write queries, you can find examples below

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
