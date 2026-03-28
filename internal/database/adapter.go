package database

import "database/sql"

// Adapter is the required contract for all types of databases
type Adapter interface {
	Connect(dsn string) error
	Close() error
	Ping() error
	GetDB() *sql.DB
	BulkInsert(tableName string, columns []string, chunk [][]any) error 
	FetchData(tableName string, limit int) (*sql.Rows, error)
}