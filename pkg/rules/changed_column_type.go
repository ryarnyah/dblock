package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type changedColumnTypeRule struct{}

func init() {
	RegisterRule("DCT01", changedColumnTypeRule{})
}

// ChangedTypeRuleError raised when a Column changed type
type ChangedTypeRuleError struct {
	RuleError

	OldColumn     string `json:"old_column_name"`
	OldColumnType string `json:"old_column_type"`
	Schema        string `json:"schema_name"`
	Table         string `json:"table_name"`
	NewColumnType string `json:"new_column_type"`
}

func (r ChangedTypeRuleError) Error() string {
	return fmt.Sprintf("[%s] colunm %s type %s of table %s.%s is incohérent in new type %s",
		r.RuleCode,
		r.OldColumn,
		r.OldColumnType,
		r.Schema,
		r.Table,
		r.NewColumnType,
	)
}

func (changedColumnTypeRule) CheckCompatibility(oldDatabase, newDatabase *model.DatabaseSchema) []error {
	return checkDatabaseSchemaTable(oldDatabase, newDatabase, func(schema model.Schema, oldTable, newTable model.TableSchema) []error {
		errors := make([]error, 0)
		for _, oldColumn := range oldTable.Columns {
			for _, newColumn := range newTable.Columns {
				if oldColumn.ColunmName == newColumn.ColunmName && oldColumn.ColumnType != newColumn.ColumnType {
					errors = append(errors, ChangedTypeRuleError{
						RuleError: RuleError{
							RuleCode: "DCT01",
						},
						OldColumn:     oldColumn.ColunmName,
						OldColumnType: oldColumn.ColumnType,
						Schema:        schema.SchemaName,
						Table:         oldTable.TableName,
						NewColumnType: newColumn.ColumnType,
					})
					break
				}
			}
		}
		return errors
	})
}
