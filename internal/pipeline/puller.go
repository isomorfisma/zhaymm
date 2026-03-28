package pipeline

import (
	"fmt"

	"github.com/isomorfisma/zhaymm/internal/database"
	"github.com/isomorfisma/zhaymm/internal/engine"
)

func RunPuller(sourceDB, targetDB database.Adapter, tableName string, maskRules map[string]string, limit int) error {
	rows, err := sourceDB.FetchData(tableName, limit)
	if err != nil {
		return fmt.Errorf("Fetch data failed %s: %w", tableName, err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	var chunk [][]any
	insertedCount := 0

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range cols {
		valuePtrs[i] = &values[i] 
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return err
		}

		rowMap := make(map[string]any)
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		if len(maskRules) > 0 {
			maskedValues, err := engine.GenerateRow(maskRules)
			if err != nil {
				return fmt.Errorf("Mask data failed: %w", err)
			}
			for k, v := range maskedValues {
				if _, exists := rowMap[k]; exists {
					rowMap[k] = v
				}
			}
		}

		var finalValues []any
		for _, col := range cols {
			finalValues = append(finalValues, rowMap[col])
		}
		chunk = append(chunk, finalValues)

		if len(chunk) == ChunkSize {
			err := targetDB.BulkInsert(tableName, cols, chunk)
			if err != nil {
				return err
			}
			insertedCount += len(chunk)
			fmt.Printf("\r-> %s: %d rows pulled & censoring...", tableName, insertedCount)
			chunk = chunk[:0]
		}
	}

	if len(chunk) > 0 {
		err := targetDB.BulkInsert(tableName, cols, chunk)
		if err != nil {
			return err
		}
		insertedCount += len(chunk)
		fmt.Printf("\r-> %s: %d rows pulled & censoring...", tableName, insertedCount)
	}

	fmt.Println("All done!")
	return nil
}