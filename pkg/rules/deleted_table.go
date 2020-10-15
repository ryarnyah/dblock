package rules

import (
	"fmt"

	"github.com/ryarnyah/dblock/pkg/model"
)

type deletedTableRule struct{}

func init() {
	RegisterRule("DT001", deletedTableRule{})
}

// DeletedTableRuleError raised when a Table has been deleted
type DeletedTableRuleError struct {
	RuleError

	Schema string `json:"schema_name"`
	Table  string `json:"table_name"`
}

func (r DeletedTableRuleError) Error() string {
	return fmt.Sprintf("[%s] table %s.%s is absent in new schema",
		r.RuleCode,
		r.Schema,
		r.Table,
	)
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
				errors = append(errors, DeletedTableRuleError{
					RuleError: RuleError{
						RuleCode: "DT001",
					},
					Schema: oldSchema.SchemaName,
					Table:  oldTable.TableName,
				})
			}
		}
		return errors
	})
}
