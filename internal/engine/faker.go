package engine

import (
	"fmt"
	"math/rand" 
	"time"     

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
		"email":       gofakeit.Email,   
		"phone":       gofakeit.Phone,  
		"company":     gofakeit.Company,	
		"word":        gofakeit.Word,    
		"random_int": func(min, max int) int {
			return gofakeit.Number(min, max)
		},
		"random_ref": func(tableName string) any {
			ids := PKStore[tableName]
			if len(ids) == 0 {
				return ""
			}
			return ids[rand.Intn(len(ids))]
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