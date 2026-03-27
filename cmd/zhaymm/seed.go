package main

import (
	"fmt"
	"log"
	"os"

	"github.com/isomorfisma/zhaymm/internal/config" 
	"github.com/isomorfisma/zhaymm/internal/database"
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