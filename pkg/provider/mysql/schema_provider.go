package mysql

import (
	"database/sql"
	"flag"
	"regexp"
	"strings"

	"github.com/ryarnyah/dblock/pkg/model"
	"github.com/ryarnyah/dblock/pkg/provider"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlSelectAllColumns = `SELECT table_schema, table_name, column_name, is_nullable, data_type, character_maximum_length, numeric_precision, numeric_scale, character_octet_length, column_default
              FROM information_schema.columns
              ORDER BY table_schema, table_name, column_name`
)

var (
	mysqlConninfo     = flag.String("mysql-conn-info", "user:password@/dbname", "MysqlQL connetion info")
	mysqlSchemaRegexp = flag.String("mysql-schema-regexp", ".*", "Regex to filter schema to process")
)

type mysqlProvider struct{}

func init() {
	provider.RegisterProvider("mysql", mysqlProvider{})
}

func (s mysqlProvider) GetCurrentModel() (*model.DatabaseSchema, error) {
	conn, err := sql.Open("mysql", *mysqlConninfo)
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

func (mysqlProvider) getModel(conn *sql.DB) ([]model.Schema, error) {
	m := make(map[string]map[string][]model.ColumnSchema)
	rows, err := conn.Query(mysqlSelectAllColumns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matcher, err := regexp.Compile(*mysqlSchemaRegexp)
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
		var numericScale string
		var characterOctetLength string

		var columnDefault string

		if err := rows.Scan(
			&schemaNameColumn,
			&tableNameColumn,
			&columnNameColumn,
			&isNullableColumn,
			&columnTypeColumn,
			&characterMaximumLengthColumn,
			&numericPrecisionColumn,
			&numericScaleColumn,
			&characterOctetLengthColumn,
			&columnDefaultColumn,
		); err != nil {
			return nil, err
		}

		if characterOctetLengthColumn.Valid {
			characterOctetLength = characterOctetLengthColumn.String
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
		if numericScaleColumn.Valid {
			numericScale = numericScaleColumn.String
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
