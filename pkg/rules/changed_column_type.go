package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type changedColumnTypeRule struct{}

func init() {
	RegisterRule("DCT01", changedColumnTypeRule{})
}

func (changedColumnTypeRule) CheckCompatibility(oldDatabase, newDatabase *model.DatabaseSchema) []error {
	return checkDatabaseSchemaTable(oldDatabase, newDatabase, func(schema model.Schema, oldTable, newTable model.TableSchema) []error {
		errors := make([]error, 0)
		for _, oldColumn := range oldTable.Columns {
			for _, newColumn := range newTable.Columns {
				if oldColumn.ColunmName == newColumn.ColunmName && oldColumn.ColumnType != newColumn.ColumnType {
					errors = append(errors, fmt.Errorf("colunm %s type %s of table %s.%s is incohérent in new type %s",
						oldColumn.ColunmName,
						oldColumn.ColumnType,
						schema.SchemaName,
						oldTable.TableName,
						newColumn.ColumnType))
					break
				}
			}
		}
		return errors
	})
}
