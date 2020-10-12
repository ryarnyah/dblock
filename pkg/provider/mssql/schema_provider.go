package mssql

import (
	"database/sql"
	"flag"
	"regexp"
	"strings"

	"github.com/ryarnyah/dblock/pkg/model"
	"github.com/ryarnyah/dblock/pkg/provider"

	// Mssql driver
	_ "github.com/denisenkom/go-mssqldb"
)

const (
	mssqlSelectAllColumns = `SELECT
                table_schema,
                table_name,
                column_name,
                is_nullable,
                data_type,
                character_maximum_length,
                numeric_precision,
                numeric_precision_radix,
                numeric_scale,
                character_octet_length,
                column_default
              FROM information_schema.columns
              ORDER BY table_schema, table_name, column_name`
)

var (
	mssqlConninfo     = flag.String("mssql-conn-info", "sqlserver://sa@localhost/SQLExpress?database=master&connection+timeout=30", "Mssql connetion info")
	mssqlSchemaRegexp = flag.String("mssql-schema-regexp", ".*", "Reex to filter schema to process")
)

type mssqlProvider struct{}

func init() {
	provider.RegisterProvider("mssql", mssqlProvider{})
}

func (s mssqlProvider) GetCurrentModel() (*model.DatabaseSchema, error) {
	conn, err := sql.Open("mssql", *mssqlConninfo)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	schemas, err := s.getModel(conn)
	if err != nil {
		return nil, err
	}
	return &model.DatabaseSchema{
		Schemas: schemas,
	}, nil
}

func (mssqlProvider) getModel(conn *sql.DB) ([]model.Schema, error) {
	m := make(map[string]map[string][]model.ColumnSchema)
	rows, err := conn.Query(mssqlSelectAllColumns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matcher, err := regexp.Compile(*mssqlSchemaRegexp)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var schemaNameColumn sql.NullString
		var tableNameColumn sql.NullString
		var columnNameColumn sql.NullString
		var isNullableColumn sql.NullString
		var columnTypeColumn sql.NullString
		var characterMaximumLengthColumn sql.NullString
		var numericPrecisionColumn sql.NullString
		var numericPrecisionRadixColumn sql.NullString
		var numericScaleColumn sql.NullString
		var characterOctetLengthColumn sql.NullString
		var columnDefaultColumn sql.NullString

		var schemaName string
		var tableName string
		var columnName string
		var isNullable string
		var columnType string
		var characterMaximumLength string
		var numericPrecision string
		var numericPrecisionRadix string
		var characterOctetLength string
		var numericScale string
		var columnDefault string

		if err := rows.Scan(
			&schemaNameColumn,
			&tableNameColumn,
			&columnNameColumn,
			&isNullableColumn,
			&columnTypeColumn,
			&characterMaximumLengthColumn,
			&numericPrecisionColumn,
			&numericPrecisionRadixColumn,
			&numericScaleColumn,
			&characterOctetLengthColumn,
			&columnDefaultColumn,
		); err != nil {
			return nil, err
		}

		if characterOctetLengthColumn.Valid {
			characterOctetLength = characterOctetLengthColumn.String
		}
		if numericScaleColumn.Valid {
			numericScale = numericScaleColumn.String
		}
		if schemaNameColumn.Valid {
			schemaName = schemaNameColumn.String
		}
		if tableNameColumn.Valid {
			tableName = tableNameColumn.String
		}
		if columnNameColumn.Valid {
			columnName = columnNameColumn.String
		}
		if isNullableColumn.Valid {
			isNullable = isNullableColumn.String
		}
		if columnTypeColumn.Valid {
			columnType = columnTypeColumn.String
		}
		if columnDefaultColumn.Valid {
			columnDefault = columnDefaultColumn.String
		}
		if characterMaximumLengthColumn.Valid {
			characterMaximumLength = characterMaximumLengthColumn.String
		}
		if numericPrecisionColumn.Valid {
			numericPrecision = numericPrecisionColumn.String
		}
		if numericPrecisionRadixColumn.Valid {
			numericPrecisionRadix = numericPrecisionRadixColumn.String
		}

		if !matcher.MatchString(schemaName) {
			continue
		}
		if _, ok := m[schemaName]; !ok {
			m[schemaName] = make(map[string][]model.ColumnSchema)
		}
		if _, ok := m[schemaName][tableName]; !ok {
			m[schemaName][tableName] = make([]model.ColumnSchema, 0)
		}
		m[schemaName][tableName] = append(m[schemaName][tableName], model.ColumnSchema{
			ColunmName: columnName,
			ColumnType: strings.Join([]string{
				columnType,
				characterMaximumLength,
				numericPrecision,
				numericPrecisionRadix,
				numericScale,
				characterOctetLength,
			}, "|"),
			NullableConstraint: isNullable == "YES",
			HasDefaultValue:    columnDefault != "",
		})
	}
	schemas := make([]model.Schema, 0)
	for schemaName, tables := range m {
		cs := model.Schema{
			SchemaName: schemaName,
			Tables:     make([]model.TableSchema, 0),
		}
		for tablename, columns := range tables {
			ct := model.TableSchema{
				TableName: tablename,
				Columns:   columns,
			}
			cs.Tables = append(cs.Tables, ct)
		}
		schemas = append(schemas, cs)
	}
	return schemas, nil
}
