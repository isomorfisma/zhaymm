package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig membaca file YAML dari hardisk dan mengembalikannya sebagai pointer ke Config
func LoadConfig(filename string) (*Config, error) {
	// 1. Baca isi file (I/O Operation)
	data, err := os.ReadFile(filename)
	if err != nil {
		// Menggunakan %w untuk membungkus (wrap) error asli agar mudah di-debug nanti
		return nil, fmt.Errorf("gagal membaca file %s: %w", filename, err)
	}	

	// 2. Siapkan 'wadah' kosong berupa Struct Config
	var cfg Config

	// 3. Terjemahkan (Unmarshal) teks YAML ke dalam memori wadah cfg
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("format YAML tidak valid: %w", err)
	}

	// 4. Kembalikan alamat memori (pointer) dari config yang sudah terisi
	return &cfg, nil
}