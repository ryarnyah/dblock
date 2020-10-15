package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type deletedColumnTableRule struct{}

func init() {
	RegisterRule("DC001", deletedColumnTableRule{})
}

// DeletedColumnRuleError raised when a Column has been deleted
type DeletedColumnRuleError struct {
	RuleError

	Column string `json:"column_name"`
	Schema string `json:"schema_name"`
	Table  string `json:"table_name"`
}

func (r DeletedColumnRuleError) Error() string {
	return fmt.Sprintf("[%s] colunm %s of table %s.%s is absent in new schema",
		r.RuleCode,
		r.Column,
		r.Schema,
		r.Table,
	)
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
				errors = append(errors, DeletedColumnRuleError{
					RuleError: RuleError{
						RuleCode: "DC001",
					},
					Column: oldColumn.ColunmName,
					Schema: schema.SchemaName,
					Table:  oldTable.TableName,
				})
			}
		}
		return errors
	})
}
