package pipeline

import (
	"fmt"


	"github.com/isomorfisma/zhaymm/internal/database"
	"github.com/isomorfisma/zhaymm/internal/engine"
)

const ChunkSize = 5000 

func RunSeeder(db database.Adapter, tableName string, columns map[string]string, totalRows int) error {
	var colNames []string
	for k := range columns {
		colNames = append(colNames, k)
	}

	var chunk [][]any
	insertedCount := 0

	for i := 0; i < totalRows; i++ {
		rowMap, err := engine.GenerateRow(columns)

		if err != nil {
			return fmt.Errorf("Failed to generate row-%d: %w", i, err)
		}
		if idVal, exists := rowMap["id"]; exists {
			engine.PKStore[tableName] = append(engine.PKStore[tableName], idVal)
		}

		var rowValues []any
		for _, colName := range colNames {
			rowValues = append(rowValues, rowMap[colName])
		}

		chunk = append(chunk, rowValues)

		if len(chunk) == ChunkSize || i == totalRows-1 {
			err := db.BulkInsert(tableName, colNames, chunk)
			if err != nil {
				return err
			}
			insertedCount += len(chunk)
			fmt.Printf("\r-> %s: %d / %d rows seeded successfully...", tableName, insertedCount, totalRows)
			
			chunk = chunk[:0] 
		}
	}

	fmt.Println("Done!")
	return nil
}