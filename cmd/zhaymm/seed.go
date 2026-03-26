package main

import (
	"fmt"
	"log"

	"github.com/isomorfisma/zhaymm/internal/config" 
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

		// Proof that Golang successfully read and understands YAML.
		fmt.Printf("Config read successfully! Found %d tables to process.\n", len(cfg.Tables))
		for _, t := range cfg.Tables {
			fmt.Printf("-> Preparing data for table '%s' with sum of %d rows.\n", t.Name, t.Count)
		}
	},
}

// Automatically runs before main()
func init() {
	// Attach seed subcommand to root
	rootCmd.AddCommand(seedCmd)
}