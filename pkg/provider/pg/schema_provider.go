package pg

import (
	"database/sql"
	"flag"
	"regexp"

	"github.com/ryarnyah/dblock/pkg/model"
	"github.com/ryarnyah/dblock/pkg/provider"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

const (
	postgresSelectAllColumns = `SELECT table_schema, table_name, column_name, is_nullable, CONCAT(data_type, '|', character_maximum_length, '|', numeric_precision, '|', numeric_precision_radix) data_type, column_default
              FROM information_schema.columns
              ORDER BY table_schema, table_name, column_name`
)

var (
	postgresConninfo     = flag.String("pg-conn-info", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=postgres", "PostgreSQL connetion info")
	postgresSchemaRegexp = flag.String("pg-schema-regexp", ".*", "Reex to filter schema to process")
)

type postgresProvider struct{}

func init() {
	provider.RegisterProvider("postgres", postgresProvider{})
}

func (s postgresProvider) GetCurrentModel() (*model.DatabaseSchema, error) {
	conn, err := sql.Open("postgres", *postgresConninfo)
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

func (postgresProvider) getModel(conn *sql.DB) ([]model.Schema, error) {
	m := make(map[string]map[string][]model.ColumnSchema)
	rows, err := conn.Query(postgresSelectAllColumns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matcher, err := regexp.Compile(*postgresSchemaRegexp)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var schemaNameColumn sql.NullString
		var tableNameColumn sql.NullString
		var columnNameColumn sql.NullString
		var isNullableColumn sql.NullString
		var columnTypeColumn sql.NullString
		var columnDefaultColumn sql.NullString

		var schemaName string
		var tableName string
		var columnName string
		var isNullable string
		var columnType string
		var columnDefault string

		if err := rows.Scan(
			&schemaNameColumn,
			&tableNameColumn,
			&columnNameColumn,
			&isNullableColumn,
			&columnTypeColumn,
			&columnDefaultColumn,
		); err != nil {
			return nil, err
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
			ColunmName:         columnName,
			ColumnType:         columnType,
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