package engine

import (
	"fmt"
	"math/rand" // Tambahkan import ini
	"time"      // Tambahkan import ini

	"github.com/expr-lang/expr"
	"github.com/brianvoe/gofakeit/v6"
)

var PKStore = make(map[string][]any)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRow(columns map[string]string) (map[string]any, error) {
	row := make(map[string]any)

	env := map[string]interface{}{
		"uuid":        gofakeit.UUID,
		"person_name": gofakeit.Name,
		"random_int": func(min, max int) int {
			return gofakeit.Number(min, max)
		},
		"random_ref": func(tableName string) any {
			ids := PKStore[tableName]
			if len(ids) == 0 {
				return "" 
			}
			randomIndex := rand.Intn(len(ids))
			return ids[randomIndex]
		},
	}

	for colName, rule := range columns {
		program, err := expr.Compile(rule, expr.Env(env))
		if err != nil {
			return nil, fmt.Errorf("aturan tidak valid pada kolom '%s': %w", colName, err)
		}

		result, err := expr.Run(program, env)
		if err != nil {
			return nil, fmt.Errorf("gagal mengeksekusi aturan pada kolom '%s': %w", colName, err)
		}

		row[colName] = result
	}

	return row, nil
}