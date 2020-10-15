package rules

import (
	"errors"
	"testing"

	"github.com/ryarnyah/dblock/pkg/model"
)

func TestChangeColumnType(t *testing.T) {
	testData := []struct {
		testName          string
		oldDatabaseSchema *model.DatabaseSchema
		newDatabaseSchema *model.DatabaseSchema
		checkErrors       []error
	}{
		{
			"change-column-type",
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
										ColumnType: "VARCHAR2",
									},
								},
							},
						},
					},
				},
			},
			[]error{
				errors.New("[DCT01] colunm ID type INTEGER of table test.test is incohérent in new type VARCHAR2"),
			},
		},
	}

	rule := changedColumnTypeRule{}
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
