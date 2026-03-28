package main

import (
	"fmt"
	"log"
	"os"

	"github.com/isomorfisma/zhaymm/internal/config"
	"github.com/isomorfisma/zhaymm/internal/dag"
	"github.com/isomorfisma/zhaymm/internal/database"
	"github.com/isomorfisma/zhaymm/internal/pipeline"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls data from the source, censors, pushes to the local target",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()

		fmt.Println("Reading censor config from schema.yaml...")
		cfg, err := config.LoadConfig("schema.yaml")
		if err != nil {
			log.Fatalf("Config error: %v", err)
		}

		graph := dag.NewGraph()
		for _, t := range cfg.Tables {
			graph.AddNode(t.Name, t.DependsOn)
		}
		executionOrder, err := graph.Sort()
		if err != nil {
			log.Fatalf("DAG error: %v", err)
		}

		fmt.Println("Opening connection to source DB and target DB...")
		sourceDSN := os.Getenv("SOURCE_DATABASE_URL")
		targetDSN := os.Getenv("DATABASE_URL")

		var sourceDB database.Adapter = &database.PostgresAdapter{}
		if err := sourceDB.Connect(sourceDSN); err != nil {
			log.Fatalf("Source DB failed to connect: %v", err)
		}
		defer sourceDB.Close()

		var targetDB database.Adapter = &database.PostgresAdapter{}
		if err := targetDB.Connect(targetDSN); err != nil {
			log.Fatalf("Target DB failed to connect: %v", err)
		}
		defer targetDB.Close()

		fmt.Println("Starting the process...")
		for _, tableName := range executionOrder {
			var maskRules map[string]string
			var limit int

			for _, t := range cfg.Tables {
				if t.Name == tableName {
					maskRules = t.Columns
					limit = t.Count 
					break
				}
			}

			err := pipeline.RunPuller(sourceDB, targetDB, tableName, maskRules, limit)
			if err != nil {
				log.Fatalf("\nError pull at table %s: %v", tableName, err)
			}
			fmt.Println()
		}

		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}