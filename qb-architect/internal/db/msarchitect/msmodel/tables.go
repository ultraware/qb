// Code generated by qb-generator; DO NOT EDIT.

package msmodel

import "git.ultraware.nl/NiseVoid/qb"

///// Tables /////
var (
	qbTablesTable = qb.Table{Name: "information_schema.tables"}

	qbTablesFTableName    = qb.TableField{Parent: &qbTablesTable, Name: "table_name"}
	qbTablesFTableSchema  = qb.TableField{Parent: &qbTablesTable, Name: "table_schema"}
	qbTablesFTableCatalog = qb.TableField{Parent: &qbTablesTable, Name: "table_catalog"}
	qbTablesFTableType    = qb.TableField{Parent: &qbTablesTable, Name: "table_type"}
)

// TablesType represents the table "Tables"
type TablesType struct {
	TableName    qb.Field
	TableSchema  qb.Field
	TableCatalog qb.Field
	TableType    qb.Field
	table        *qb.Table
}

// GetTable returns an object with info about the table
func (t *TablesType) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *TablesType) Select(f ...qb.Field) *qb.SelectBuilder {
	return t.table.Select(f)
}

// Delete creates a DELETE query
func (t *TablesType) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts an UPDATE query
func (t *TablesType) Update() *qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts an INSERT query
func (t *TablesType) Insert(f ...qb.Field) *qb.InsertBuilder {
	return t.table.Insert(f)
}

// Tables returns a new TablesType
func Tables() *TablesType {
	table := qbTablesTable
	return &TablesType{
		qbTablesFTableName.Copy(&table),
		qbTablesFTableSchema.Copy(&table),
		qbTablesFTableCatalog.Copy(&table),
		qbTablesFTableType.Copy(&table),
		&table,
	}
}

///// Columns /////
var (
	qbColumnsTable = qb.Table{Name: "information_schema.columns"}

	qbColumnsFColumnName             = qb.TableField{Parent: &qbColumnsTable, Name: "column_name"}
	qbColumnsFTableSchema            = qb.TableField{Parent: &qbColumnsTable, Name: "table_schema"}
	qbColumnsFDataType               = qb.TableField{Parent: &qbColumnsTable, Name: "data_type"}
	qbColumnsFIsNullable             = qb.TableField{Parent: &qbColumnsTable, Name: "is_nullable"}
	qbColumnsFCharacterMaximumLength = qb.TableField{Parent: &qbColumnsTable, Name: "character_maximum_length"}
	qbColumnsFTableCatalog           = qb.TableField{Parent: &qbColumnsTable, Name: "table_catalog"}
	qbColumnsFTableName              = qb.TableField{Parent: &qbColumnsTable, Name: "table_name"}
)

// ColumnsType represents the table "Columns"
type ColumnsType struct {
	ColumnName             qb.Field
	TableSchema            qb.Field
	DataType               qb.Field
	IsNullable             qb.Field
	CharacterMaximumLength qb.Field
	TableCatalog           qb.Field
	TableName              qb.Field
	table                  *qb.Table
}

// GetTable returns an object with info about the table
func (t *ColumnsType) GetTable() *qb.Table {
	return t.table
}

// Select starts a SELECT query
func (t *ColumnsType) Select(f ...qb.Field) *qb.SelectBuilder {
	return t.table.Select(f)
}

// Delete creates a DELETE query
func (t *ColumnsType) Delete(c1 qb.Condition, c ...qb.Condition) qb.Query {
	return t.table.Delete(c1, c...)
}

// Update starts an UPDATE query
func (t *ColumnsType) Update() *qb.UpdateBuilder {
	return t.table.Update()
}

// Insert starts an INSERT query
func (t *ColumnsType) Insert(f ...qb.Field) *qb.InsertBuilder {
	return t.table.Insert(f)
}

// Columns returns a new ColumnsType
func Columns() *ColumnsType {
	table := qbColumnsTable
	return &ColumnsType{
		qbColumnsFColumnName.Copy(&table),
		qbColumnsFTableSchema.Copy(&table),
		qbColumnsFDataType.Copy(&table),
		qbColumnsFIsNullable.Copy(&table),
		qbColumnsFCharacterMaximumLength.Copy(&table),
		qbColumnsFTableCatalog.Copy(&table),
		qbColumnsFTableName.Copy(&table),
		&table,
	}
}