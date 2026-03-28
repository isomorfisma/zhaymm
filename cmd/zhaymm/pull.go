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

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Extract data from source, anonymize, and load into target database",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		sourceURL, _ := cmd.Flags().GetString("source-db")
		targetURL, _ := cmd.Flags().GetString("target-db")

		fmt.Printf("[1/4] Reading masking rules from %s...\n", configPath)
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Fatal config error: %v", err)
		}

		fmt.Println("[2/4] Analyzing table relations (DAG)...")
		graph := dag.NewGraph()
		for _, t := range cfg.Tables {
			graph.AddNode(t.Name, t.DependsOn)
		}
		executionOrder, err := graph.Sort()
		if err != nil {
			log.Fatalf("DAG Error: %v", err)
		}

		fmt.Println("[3/4] Connecting to Source and Target databases...")
		var sourceDB database.Adapter = &database.PostgresAdapter{}
		if err := sourceDB.Connect(sourceURL); err != nil {
			log.Fatalf("Source DB connection failed: %v", err)
		}
		defer sourceDB.Close()

		var targetDB database.Adapter = &database.PostgresAdapter{}
		if err := targetDB.Connect(targetURL); err != nil {
			log.Fatalf("Target DB connection failed: %v", err)
		}
		defer targetDB.Close()
		fmt.Println("      -> Connections established.")

		fmt.Println("[4/4] Starting pull and anonymization pipeline...")
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
				log.Fatalf("\nError pulling table %s: %v", tableName, err)
			}
			fmt.Println()
		}

		fmt.Println("Pull & anonymize completed!")
	},
}

func init() {
	pullCmd.Flags().StringP("config", "c", "schema.yaml", "Path to the YAML schema file")
	pullCmd.Flags().StringP("source-db", "s", "", "Source database connection string (DSN)")
	pullCmd.Flags().StringP("target-db", "t", "", "Target database connection string (DSN)")
	
	pullCmd.MarkFlagRequired("source-db")
	pullCmd.MarkFlagRequired("target-db")

	rootCmd.AddCommand(pullCmd)
}