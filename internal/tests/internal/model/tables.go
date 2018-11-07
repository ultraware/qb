package model

import "git.ultraware.nl/NiseVoid/qb"

///// One /////
var qbOneTable = qb.Table{Name: "one"}

var qbOneFID = qb.TableField{Parent: &qbOneTable, Name: "ID", ReadOnly: true, Type: qb.Int, Size: 32}
var qbOneFName = qb.TableField{Parent: &qbOneTable, Name: "Name", Type: qb.String, Size: 50}
var qbOneFCreatedAt = qb.TableField{Parent: &qbOneTable, Name: "CreatedAt", ReadOnly: true, Type: qb.Time}

// OneType represents the table "One"
type OneType struct {
	ID        qb.Field
	Name      qb.Field
	CreatedAt qb.Field
	table     *qb.Table
}

// All returns every field as an array
func (t *OneType) All() []qb.Field {
	return []qb.Field{t.ID, t.Name, t.CreatedAt}
}

// GetTable returns an object with info about the table
func (t *OneType) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *OneType) Select(f ...qb.Field) *qb.SelectBuilder {
	return t.table.Select(f)
}

// Delete creates a DELETE query
func (t *OneType) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts an UPDATE query
func (t *OneType) Update() *qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts an INSERT query
func (t *OneType) Insert(f ...qb.Field) *qb.InsertBuilder {
	return t.table.Insert(f)
}

// One returns a new OneType
func One() *OneType {
	table := qbOneTable
	return &OneType{
		qbOneFID.Copy(&table),
		qbOneFName.Copy(&table),
		qbOneFCreatedAt.Copy(&table),
		&table,
	}
}

///// Two /////
var qbTwoTable = qb.Table{Name: "two", Alias: "tw"}

var qbTwoFOneID = qb.TableField{Parent: &qbTwoTable, Name: "OneID", Type: qb.Int, Size: 32}
var qbTwoFNumber = qb.TableField{Parent: &qbTwoTable, Name: "Number", Type: qb.Int, Size: 32}
var qbTwoFComment = qb.TableField{Parent: &qbTwoTable, Name: "Comment", Type: qb.String, Size: 100}
var qbTwoFModifiedAt = qb.TableField{Parent: &qbTwoTable, Name: "ModifiedAt", Type: qb.Time, Nullable: true}

// TwoType represents the table "Two"
type TwoType struct {
	OneID      qb.Field
	Number     qb.Field
	Comment    qb.Field
	ModifiedAt qb.Field
	table      *qb.Table
}

// All returns every field as an array
func (t *TwoType) All() []qb.Field {
	return []qb.Field{t.OneID, t.Number, t.Comment, t.ModifiedAt}
}

// GetTable returns an object with info about the table
func (t *TwoType) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *TwoType) Select(f ...qb.Field) *qb.SelectBuilder {
	return t.table.Select(f)
}

// Delete creates a DELETE query
func (t *TwoType) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts an UPDATE query
func (t *TwoType) Update() *qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts an INSERT query
func (t *TwoType) Insert(f ...qb.Field) *qb.InsertBuilder {
	return t.table.Insert(f)
}

// Two returns a new TwoType
func Two() *TwoType {
	table := qbTwoTable
	return &TwoType{
		qbTwoFOneID.Copy(&table),
		qbTwoFNumber.Copy(&table),
		qbTwoFComment.Copy(&table),
		qbTwoFModifiedAt.Copy(&table),
		&table,
	}
}
