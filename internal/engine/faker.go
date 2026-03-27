package engine

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/brianvoe/gofakeit/v6"
)

// GenerateRow creates one row of map[string]any data depends on the column
func GenerateRow(columns map[string]string) (map[string]any, error) {
	row := make(map[string]any)

	// Environment for Expression Engine
	env := map[string]interface{}{
		"uuid":        gofakeit.UUID,
		"person_name": gofakeit.Name,
		
		"random_int": func(min, max int) int {
			return gofakeit.Number(min, max)
		},
	}

	// 2. Loop every column from the YAML
	for colName, rule := range columns {
		
		program, err := expr.Compile(rule, expr.Env(env))
		if err != nil {
			return nil, fmt.Errorf("Invalid rule at column: '%s': %w", colName, err)
		}

		result, err := expr.Run(program, env)
		if err != nil {
			return nil, fmt.Errorf("Failed to execute rule at column: '%s': %w", colName, err)
		}

		row[colName] = result
	}

	return row, nil
}