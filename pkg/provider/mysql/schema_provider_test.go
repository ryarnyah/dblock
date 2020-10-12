package mysql

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryarnyah/dblock/pkg/model"
)

func TestLoadDatabaseSchemaFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{
		"table_schema",
		"table_name",
		"column_name",
		"is_nullable",
		"data_type",
		"character_maximum_length",
		"numeric_precision",
		"numeric_scale",
		"character_octet_length",
		"column_default",
	}).AddRow(
		"test_schema",
		"test_table",
		"test_column",
		"NO",
		"VARCHAR2",
		"12",
		nil,
		nil,
		nil,
		"",
	).AddRow(
		"test_schema",
		"test_table",
		"second_column",
		"YES",
		"INTEGER",
		nil,
		"12",
		"2",
		nil,
		"12",
	)
	mock.ExpectQuery(".*").
		WillReturnRows(rows)

	expectedSchemas := []model.Schema{
		{
			SchemaName: "test_schema",
			Tables: []model.TableSchema{
				{
					TableName: "test_table",
					Columns: []model.ColumnSchema{
						{
							ColunmName:         "test_column",
							ColumnType:         "VARCHAR2|12|||",
							NullableConstraint: false,
							HasDefaultValue:    false,
						},
						{
							ColunmName:         "second_column",
							ColumnType:         "INTEGER||12|2|",
							NullableConstraint: true,
							HasDefaultValue:    true,
						},
					},
				},
			},
		},
	}

	pgProvider := mysqlProvider{}
	schemas, err := pgProvider.getModel(db)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedSchemas, schemas) {
		t.Fatalf("got %+v, expectd %+v", schemas, expectedSchemas)
	}
}
