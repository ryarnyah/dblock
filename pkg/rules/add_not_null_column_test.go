package rules

import (
	"errors"
	"testing"

	"github.com/ryarnyah/dblock/pkg/model"
)

func TestAddNewColumn(t *testing.T) {
	testData := []struct {
		testName          string
		oldDatabaseSchema *model.DatabaseSchema
		newDatabaseSchema *model.DatabaseSchema
		checkErrors       []error
	}{
		{
			"no-error",
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
								},
							},
						},
					},
				},
			},
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
								},
							},
						},
					},
				},
			},
			[]error{},
		},
		{
			"add-new-column-default-value",
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
								},
							},
						},
					},
				},
			},
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
									{
										ColunmName:         "ID2",
										ColumnType:         "INTEGER",
										NullableConstraint: true,
										HasDefaultValue:    true,
									},
								},
							},
						},
					},
				},
			},
			[]error{},
		},
		{
			"add-column-without-default-value",
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
								},
							},
						},
					},
				},
			},
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
									{
										ColunmName:         "ID2",
										ColumnType:         "INTEGER",
										NullableConstraint: true,
										HasDefaultValue:    false,
									},
								},
							},
						},
					},
				},
			},
			[]error{},
		},
		{
			"add-column-without-default-value-not-null",
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
								},
							},
						},
					},
				},
			},
			&model.DatabaseSchema{
				Schemas: []model.Schema{
					{
						SchemaName: "test",
						Tables: []model.TableSchema{
							{
								TableName: "test",
								Columns: []model.ColumnSchema{
									{
										ColunmName: "ID",
										ColumnType: "INTEGER",
									},
									{
										ColunmName:         "ID2",
										ColumnType:         "INTEGER",
										NullableConstraint: false,
										HasDefaultValue:    false,
									},
								},
							},
						},
					},
				},
			},
			[]error{
				errors.New("[DCT02] not null column ID2 of test.test added without default value"),
			},
		},
	}

	rule := addColumnNotNullTypeRule{}
	for _, tt := range testData {
		t.Run(tt.testName, func(t *testing.T) {
			errs := rule.CheckCompatibility(tt.oldDatabaseSchema, tt.newDatabaseSchema)
			if len(errs) != len(tt.checkErrors) {
				t.Errorf("got %+v, expect %+v", len(errs), len(tt.checkErrors))
			}
			for i, err := range errs {
				if err.Error() != tt.checkErrors[i].Error() {
					t.Errorf("got '%+v', expect '%+v'", err, tt.checkErrors[i])
				}
			}
		})
	}
}
