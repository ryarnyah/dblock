package rules

import (
	"github.com/ryarnyah/dblock/pkg/model"
)

// RegistredRules rules available
var RegistredRules = make(map[string]Rule)

// Rule represent a DB Rule
type Rule interface {
	CheckCompatibility(*model.DatabaseSchema, *model.DatabaseSchema) []error
}

// RegisterRule register rule to be processed
func RegisterRule(code string, rule Rule) {
	RegistredRules[code] = rule
}

// RuleError code + error
type RuleError struct {
	RuleCode string `json:"rule_code"`
}

func checkDatabaseSchemaTable(oldDatabase, newDatabase *model.DatabaseSchema, checkFunction func(model.Schema, model.TableSchema, model.TableSchema) []error) []error {
	return checkDatabaseSchemaSchema(oldDatabase, newDatabase, func(oldSchema, newSchema model.Schema) []error {
		errors := make([]error, 0)
		for _, oldTable := range oldSchema.Tables {
			for _, newTable := range newSchema.Tables {
				if oldTable.TableName == newTable.TableName {
					errors = append(errors, checkFunction(oldSchema, oldTable, newTable)...)
					break
				}
			}
		}
		return errors
	})
}

func checkDatabaseSchemaSchema(oldDatabase, newDatabase *model.DatabaseSchema, checkFunction func(model.Schema, model.Schema) []error) []error {
	errors := make([]error, 0)
	for _, oldSchema := range oldDatabase.Schemas {
		for _, newSchema := range newDatabase.Schemas {
			if oldSchema.SchemaName == newSchema.SchemaName {
				errors = append(errors, checkFunction(oldSchema, newSchema)...)
				break
			}
		}
	}
	return errors
}
