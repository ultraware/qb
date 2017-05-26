package qbdb

import "git.ultraware.nl/NiseVoid/qb"

// BatchInsert is used to insert multiple records at once
type BatchInsert struct {
	count    int
	sql      string
	Values   []interface{}
	upsert   bool
	conflict []qb.DataField
	fields   []qb.DataField
}

// NewBatch returns a BatchInsert
func (db QueryTarget) NewBatch(record Savable) *BatchInsert {
	f := record.All()
	s := qb.InsertHeaderSQL(record.GetTable(), f)
	return &BatchInsert{sql: s, fields: f}
}

// Upsert makes the batch an upsert batch instead of just a plain insert
// Not all Drivers support this
func (b *BatchInsert) Upsert(c ...qb.DataField) {
	b.upsert = true
	b.conflict = c
}

// Add adds the record to the batch
func (b *BatchInsert) Add(record Savable) {
	s, v := qb.InsertValueSQL(record.All())
	if b.count > 0 {
		b.sql += `,` + "\n"
	}
	b.sql += s
	b.Values = append(b.Values, v...)
	b.count++
}

// ExecuteBatch executes the given batch
func (db QueryTarget) ExecuteBatch(b *BatchInsert) error {
	if b.count == 0 {
		panic(`Cannot insert empty batch`)
	}
	s, v := b.sql+"\n", b.Values
	if len(b.conflict) > 0 {
		s += db.Driver.UpsertSQL(b.fields, b.conflict)
	}

	return db.prepareExec(s, v)
}
