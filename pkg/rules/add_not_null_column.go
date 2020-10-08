package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type addColumnNotNullTypeRule struct{}

func init() {
	RegisterRule("DCT02", addColumnNotNullTypeRule{})
}

func (addColumnNotNullTypeRule) CheckCompatibility(oldDatabase, newDatabase *model.DatabaseSchema) []error {
	return checkDatabaseSchemaTable(oldDatabase, newDatabase, func(schema model.Schema, oldTable, newTable model.TableSchema) []error {
		errors := make([]error, 0)
		for _, newColumn := range newTable.Columns {
			isNewColunm := true
			for _, oldColumn := range oldTable.Columns {
				if oldColumn.ColunmName == newColumn.ColunmName {
					isNewColunm = false
				}
			}
			if isNewColunm && !newColumn.NullableConstraint && !newColumn.HasDefaultValue {
				errors = append(errors, fmt.Errorf("not null column %s of %s.%s added without default value",
					newColumn.ColunmName,
					schema.SchemaName,
					oldTable.TableName))
				break
			}
		}
		return errors
	})
}
