package model

// ColumnSchema DB column representation
type ColumnSchema struct {
	ColunmName         string `json:"column_name"`
	ColumnType         string `json:"column_type"`
	NullableConstraint bool   `json:"nullable_constraint"`
	HasDefaultValue    bool   `json:"has_default_value"`
}

// TableSchema DB table representation
type TableSchema struct {
	TableName string         `json:"table_name"`
	Columns   []ColumnSchema `json:"table_columns"`
}

// Schema DB schema representation
type Schema struct {
	SchemaName string        `json:"schema_name"`
	Tables     []TableSchema `json:"schema_tables"`
}

// DatabaseSchema contain all database schemas
type DatabaseSchema struct {
	Schemas []Schema `json:"database_schemas"`
}
