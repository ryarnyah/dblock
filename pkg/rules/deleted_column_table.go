package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type deletedColumnTableRule struct{}

func init() {
	RegisterRule("DC001", deletedColumnTableRule{})
}

func (deletedColumnTableRule) CheckCompatibility(oldDatabase, newDatabase *model.DatabaseSchema) []error {
	return checkDatabaseSchemaTable(oldDatabase, newDatabase, func(schema model.Schema, oldTable, newTable model.TableSchema) []error {
		errors := make([]error, 0)
		for _, oldColumn := range oldTable.Columns {
			colunmExist := false
			for _, newColumn := range newTable.Columns {
				if oldColumn.ColunmName == newColumn.ColunmName {
					colunmExist = true
					break
				}
			}
			if !colunmExist {
				errors = append(errors, fmt.Errorf("colunm %s of table %s.%s is absent in new schema",
					oldColumn.ColunmName,
					schema.SchemaName,
					oldTable.TableName))
			}
		}
		return errors
	})
}
