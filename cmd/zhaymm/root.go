package main

import (
	"github.com/spf13/cobra"
)

// rootCmd mewakili perintah dasar ketika dipanggil tanpa sub-perintah
var rootCmd = &cobra.Command{
	Use:   "zhaymm",
	Short: "zhaymm is a blazing-fast database seeder",
	Long:  `A zero-memory-bloat database seeder & anonymization pipeline powered by Go.`,
	// Ini yang dieksekusi kalau user cuma ketik 'zhaymm' tanpa embel-embel
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute menambahkan semua child commands ke root command dan men-set flags.
func Execute() error {
	return rootCmd.Execute()
}