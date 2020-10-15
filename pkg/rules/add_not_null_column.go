package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type addColumnNotNullTypeRule struct{}

func init() {
	RegisterRule("DCT02", addColumnNotNullTypeRule{})
}

// NotNullRuleError raised when a not null Column is added
type NotNullRuleError struct {
	RuleError

	Column string `json:"column_name"`
	Schema string `json:"schema_name"`
	Table  string `json:"table_name"`
}

func (r NotNullRuleError) Error() string {
	return fmt.Sprintf("[%s] not null column %s of %s.%s added without default value",
		r.RuleCode,
		r.Column,
		r.Schema,
		r.Table,
	)
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
				errors = append(errors, NotNullRuleError{
					RuleError: RuleError{
						RuleCode: "DCT02",
					},
					Column: newColumn.ColunmName,
					Schema: schema.SchemaName,
					Table:  oldTable.TableName,
				})
				break
			}
		}
		return errors
	})
}
