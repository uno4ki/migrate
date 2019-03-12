package migrate

import (
	"database/sql"
	"fmt"
	"log"
)

// Migrator seed/migrate callback
type migrator_t func(*sql.DB) error

// Schema/Seed migrations tables
const (
	schema_migrations = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(128) NOT NULL,
			PRIMARY KEY(version)
		)`

	seed_migrations = `
		CREATE TABLE IF NOT EXISTS seed_migrations (
			version VARCHAR(128) NOT NULL,
			PRIMARY KEY(version)
		)`
)

// Create migrations/seeds table if not exists
func createIfNot(db *sql.DB) error {
	if _, err := db.Exec(schema_migrations); err != nil {
		return err
	}

	if _, err := db.Exec(seed_migrations); err != nil {
		return err
	}

	return nil
}

// Get list of currently applied migrations.
func applied(db *sql.DB, table string) []string {
	var applied []string

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s_migrations", table))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		applied = append(applied, name)
	}

	return applied
}

// Run registered migrations
func migrate(db *sql.DB) error {
	if err := createIfNot(db); err != nil {
		return err
	}

	current := applied(db, "schema")

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, m := range Migrations.select_migrations(current) {
		log.Printf("=== Migrate: %s\n", m.name)

		if err := m.execute(tx); err != nil {
			log.Printf("!!! Failed to migrate: %s\n", err)
			tx.Rollback()
			return err
		}

		stmt, err := tx.Prepare("INSERT INTO schema_migrations (version)VALUES($1)")
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmt.Exec(m.name)
		if err != nil {
			tx.Rollback()
			return err
		}
		stmt.Close()
	}

	return tx.Commit()
}

// Run registered seeds
func seed(db *sql.DB) error {
	if err := createIfNot(db); err != nil {
		return err
	}

	current := applied(db, "seed")

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, m := range Seeds.select_migrations(current) {
		log.Printf("=== Seed: %s\n", m.name)

		if err := m.execute(tx); err != nil {
			log.Printf("!!! Failed to seed: %s\n", err)
			tx.Rollback()
			return err
		}

		stmt, err := tx.Prepare("INSERT INTO seed_migrations (version) VALUES ($1)")
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmt.Exec(m.name)
		if err != nil {
			tx.Rollback()
			return err
		}
		stmt.Close()
	}

	return tx.Commit()
}
