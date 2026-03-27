package main

import (
	"fmt"
	"log"
	"os"

	"github.com/isomorfisma/zhaymm/internal/config" 
	"github.com/isomorfisma/zhaymm/internal/database"
	"github.com/isomorfisma/zhaymm/internal/dag"
	"github.com/isomorfisma/zhaymm/internal/engine"
	"github.com/joho/godotenv" 
	"github.com/spf13/cobra"
)

// seedCmd is for the CLI 'seed command'
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Starting seeding process to database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reading schema.yaml configuration...")
		
		// Calls LoadConfig function inside parser.go
		cfg, err := config.LoadConfig("schema.yaml")
		if err != nil {
			log.Fatalf("Fatal error: %v", err)
		}

		fmt.Println("Analyzing table relations...")
		graph := dag.NewGraph()

		// Inserts all table to the graph
		for _, t := range cfg.Tables {
			graph.AddNode(t.Name, t.DependsOn)
		}

		// Run the sorting algorithm
		executionOrder, err := graph.Sort()
		if err != nil {
			log.Fatalf("DAG Error: %v", err)
		}

		fmt.Printf("-> Safe execution order (from left to right): %v\n", executionOrder)

		fmt.Println("Testing data factory (generating 1 row of data)...")
		firstTable := executionOrder[0]
		var targetCols map[string]string
		for _, t := range cfg.Tables{
			if t.Name == firstTable {
				targetCols = t.Columns
				break
			}
		}

		dummyRow, err := engine.GenerateRow(targetCols)
		if err != nil {
			log.Fatalf("Failed to generate data: %v", err)
		}

		fmt.Printf("Data generated successfully for table '%s':\n", firstTable)
		for col, val := range dummyRow {
			fmt.Printf("	- %s %v\n",col, val)
		}

		fmt.Printf("-> Found %d table.\n", len(cfg.Tables))
		fmt.Println("Trying to connect to database...")

		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("Error: DATABASE_URL not set")
		}
		var dbAdapter database.Adapter = &database.PostgresAdapter{}

		err = dbAdapter.Connect(dsn)
		if err!= nil {
			log.Fatalf("Failed to connect to database. %v\nMake sure Postgres is running!", err)
		}
		defer dbAdapter.Close()
		fmt.Println("-> Successfully connected to database.")
	},

}

// Automatically runs before main()
func init() {
	// Attach seed subcommand to root
	rootCmd.AddCommand(seedCmd)
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment OS bawaan.")
	}
}