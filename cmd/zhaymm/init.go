package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a boilerplate schema.yaml configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the target filename from flags
		filename, _ := cmd.Flags().GetString("file")

		// Check if file already exists to prevent accidental overwrites
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("Warning: File '%s' already exists.\n", filename)
			fmt.Println("Aborting operation to prevent overwriting your existing configuration.")
			return
		}

		// The YAML Template Boilerplate
		template := `# ==========================================
# Database Seeder Configuration
# ==========================================
# Documentation: https://github.com/isomorfisma/zhaymm

tables:
  # Example 1: An independent table
  - name: users
    count: 10
    depends_on: []
    columns:
      id: "uuid()"
      name: "person_name()"
      email: "email()"
      
  # Example 2: A dependent table with Foreign Key (Uncomment to use)
  # - name: orders
  #   count: 50
  #   depends_on: ["users"]
  #   columns:
  #     id: "uuid()"
  #     user_id: "random_ref('users')"
  #     total: "random_int(10000, 500000)"
`

		// Write the template to the file
		err := os.WriteFile(filename, []byte(template), 0644)
		if err != nil {
			fmt.Printf("Failed to create file: %v\n", err)
			return
		}

		fmt.Printf("Successfully generated boilerplate at '%s'!\n", filename)
		fmt.Println("   -> Next step: Edit the file to match your database, then run './zhaymm seed --db <YOUR_DB_URL>'")
	},
}

func init() {
	// Add a flag just in case the user wants to name it something other than schema.yaml
	initCmd.Flags().StringP("file", "f", "schema.yaml", "Name of the generated config file")
	
	rootCmd.AddCommand(initCmd)
}