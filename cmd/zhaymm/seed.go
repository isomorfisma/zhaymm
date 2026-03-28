package main

import (
	"fmt"
	"log"

	"github.com/isomorfisma/zhaymm/internal/config"
	"github.com/isomorfisma/zhaymm/internal/dag"
	"github.com/isomorfisma/zhaymm/internal/database"
	"github.com/isomorfisma/zhaymm/internal/pipeline"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Start seeding mock data into the database",
	Run: func(cmd *cobra.Command, args []string) {
		// Get values from CLI flags
		configPath, _ := cmd.Flags().GetString("config")
		dbURL, _ := cmd.Flags().GetString("db")

		// Load YAML Configuration
		fmt.Printf("[1/4] Reading configuration from %s...\n", configPath)
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Fatal config error: %v", err)
		}

		// Analyze DAG (Directed Acyclic Graph)
		fmt.Println("[2/4] Analyzing table relations (DAG)...")
		graph := dag.NewGraph()
		for _, t := range cfg.Tables {
			graph.AddNode(t.Name, t.DependsOn)
		}
		executionOrder, err := graph.Sort()
		if err != nil {
			log.Fatalf("DAG Error: %v", err)
		}
		fmt.Printf("      -> Safe execution order: %v\n", executionOrder)

		// Connect to Database
		fmt.Println("[3/4] Connecting to target database...")
		var dbAdapter database.Adapter = &database.PostgresAdapter{}
		err = dbAdapter.Connect(dbURL)
		if err != nil {
			log.Fatalf("Database connection failed: %v", err)
		}
		defer dbAdapter.Close()
		fmt.Println("      -> Connection established.")

		// Execute Pipeline
		fmt.Println("[4/4] Starting data seeding pipeline...")
		for _, tableName := range executionOrder {
			var targetTable *config.Table
			for _, t := range cfg.Tables {
				if t.Name == tableName {
					targetTable = &t
					break
				}
			}

			err := pipeline.RunSeeder(dbAdapter, targetTable.Name, targetTable.Columns, targetTable.Count)
			if err != nil {
				log.Fatalf("\nError seeding table %s: %v", targetTable.Name, err)
			}
			fmt.Println() 
		}

		fmt.Println("All tables seeded successfully.")
	},
}

func init() {
	// Define flags for the 'seed' command
	// StringP takes: name, shorthand, default value, description
	seedCmd.Flags().StringP("config", "c", "schema.yaml", "Path to the YAML schema file")
	seedCmd.Flags().StringP("db", "d", "", "Target database connection string (DSN)")
	
	// Mark 'db' as a required flag. The CLI will reject the command if it's missing.
	seedCmd.MarkFlagRequired("db")

	rootCmd.AddCommand(seedCmd)
}