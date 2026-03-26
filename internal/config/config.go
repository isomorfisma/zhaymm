package config

// Config is the main representation from schema.yaml contents
type Config struct {
	Tables []Table `yaml:"tables"`
}

// Table represent configuration for 1 specific table
type Table struct {
	Name    string            `yaml:"name"`
	Count   int               `yaml:"count"`
	Columns map[string]string `yaml:"columns"`
}

