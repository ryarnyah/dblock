package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type deletedTableRule struct{}

func init() {
	RegisterRule("DT001", deletedTableRule{})
}

func (deletedTableRule) CheckCompatibility(oldDatabase, newDatabase *model.DatabaseSchema) []error {
	return checkDatabaseSchemaSchema(oldDatabase, newDatabase, func(oldSchema, newSchema model.Schema) []error {
		errors := make([]error, 0)
		for _, oldTable := range oldSchema.Tables {
			tableExist := false
			for _, newTable := range newSchema.Tables {
				if oldTable.TableName == newTable.TableName {
					tableExist = true
					break
				}
			}
			if !tableExist {
				errors = append(errors, fmt.Errorf("table %s.%s is absent in new schema",
					oldSchema.SchemaName,
					oldTable.TableName))
			}
		}
		return errors
	})
}
