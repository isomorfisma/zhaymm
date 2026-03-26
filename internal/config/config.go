package config

// Config adalah representasi utama dari isi schema.yaml
type Config struct {
	Tables []Table `yaml:"tables"`
}

// Table mewakili konfigurasi untuk satu tabel spesifik
type Table struct {
	Name    string            `yaml:"name"`
	Count   int               `yaml:"count"`
	Columns map[string]string `yaml:"columns"`
}
