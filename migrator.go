package migrate

import (
	"database/sql"
	"fmt"
)

// Migrator
type Migrator struct {
	driver   string // Driver name
	dsn      string // DSN for seed/migrate operations
	root_dsn string // DSN for create/drop database operatoins
}

// Initialize new migrator instance.
func NewMigrator(driver, dsn, root_dsn string) *Migrator {
	return &Migrator{driver, dsn, root_dsn}
}

// Connect to the given dsn
func connect(driver, dsn string) (*sql.DB, error) {
	var err error
	var db *sql.DB

	if db, err = sql.Open(driver, dsn); err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Execute
func execute(driver, dsn string, exec migrator_t, cmd string) error {
	var db *sql.DB

	db, err := connect(driver, dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	if exec != nil {
		return exec(db)
	} else {
		_, err := db.Exec(cmd)
		return err
	}
}

// Create database
func (m *Migrator) Create(dbname string) error {
	return execute(m.driver, m.root_dsn, nil, fmt.Sprintf("CREATE DATABASE %s", dbname))
}

// Drop database
func (m *Migrator) Drop(dbname string) error {
	return execute(m.driver, m.root_dsn, nil, fmt.Sprintf("DROP DATABASE %s", dbname))
}

// Execute migrations
func (m *Migrator) Migrate() error {
	return execute(m.driver, m.dsn, migrate, "")
}

// Execute seed
func (m *Migrator) Seed() error {
	return execute(m.driver, m.dsn, seed, "")
}
