package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Implementation of database interface
type PostgresAdapter struct {
	db *sql.DB
}

// Connect opens connection to PostgreSQL and manages the Connection Pool
func (p *PostgresAdapter) Connect(dsn string) error {
	var err error
	p.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("gagal membuka koneksi awal: %w", err)
	}

	// 20 max open connections silmutaneously
	p.db.SetMaxOpenConns(20)
	// 5 max idle connections
	p.db.SetMaxIdleConns(5)
	// Closing connections unused for 15 minutes
	p.db.SetConnMaxLifetime(15 * time.Minute)

	return p.Ping()
}

func (p *PostgresAdapter) Ping() error {
	if err := p.db.Ping(); err != nil {
		return fmt.Errorf("Fail to ping database: %w", err)
	}
	return nil
}

func (p *PostgresAdapter) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresAdapter) GetDB() *sql.DB {
	return p.db
}