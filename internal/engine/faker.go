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

func GenerateRow(maskRules map[string]string) (map[string]any, error) {
	result := make(map[string]any)

	// The Arsenal: Registering available faker functions for the user
	env := map[string]interface{}{
		"uuid":        gofakeit.UUID,
		"person_name": gofakeit.Name,
		"first_name":  gofakeit.FirstName,
		"last_name":   gofakeit.LastName,
		"gender":      gofakeit.Gender,
		"ssn":         gofakeit.SSN,
		"job_title":   gofakeit.JobTitle,

		"email":       gofakeit.Email,
		"phone":       gofakeit.Phone,
		"username":    gofakeit.Username,
		"password":    func() string { return gofakeit.Password(true, true, true, false, false, 12) },
		"ipv4":        gofakeit.IPv4Address,
		"mac_address": gofakeit.MacAddress,
		"url":         gofakeit.URL,

		"city":      gofakeit.City,
		"country":   gofakeit.Country,
		"street":    gofakeit.Street,
		"zip":       gofakeit.Zip,
		"latitude":  gofakeit.Latitude,
		"longitude": gofakeit.Longitude,

		"company":     gofakeit.Company,
		"credit_card": func() string { return gofakeit.CreditCard().Number },
		"currency":    gofakeit.CurrencyShort,
		"price":       gofakeit.Price,

		"word":      gofakeit.Word,
		"sentence":  gofakeit.Sentence,
		"paragraph": gofakeit.Paragraph,
		"color":     gofakeit.Color,

		"date": gofakeit.Date,
		"year": gofakeit.Year,

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

	for colName, expression := range maskRules {
		// Compile the string expression (e.g., "person_name()") into executable code
		program, err := expr.Compile(expression, expr.Env(env))
		if err != nil {
			return nil, fmt.Errorf("failed to compile expression for column '%s': %w", colName, err)
		}

		// Execute the compiled program
		output, err := expr.Run(program, env)
		if err != nil {
			return nil, fmt.Errorf("failed to execute expression for column '%s': %w", colName, err)
		}

		result[colName] = output
	}

	return result, nil
}