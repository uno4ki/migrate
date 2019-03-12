// Sample MIGRATE file
package main

import (
	"database/sql"
	m "github.com/uno4ki/migrate"
)

// Simple migrations flow.
func init() {

	// Create extension
	m.Migrations.Register(1, "int array extension", func(tx *sql.Tx) error {
		return m.ExecStmt(tx, "CREATE EXTENSION intarray;")
	})

	// Create users table
	m.Migrations.Register(2, "create users table", func(tx *sql.Tx) error {
		return m.ExecStmt(tx, `
			CREATE TABLE users (
				id              BIGSERIAL PRIMARY KEY,
				name            VARCHAR(128) NOT NULL
			);
		`)
	})

	// Create email field + unique index on email
	m.Migrations.Register(3, "add email to users", func(tx *sql.Tx) error {
		if err := m.ExecStmt(tx, `ALTER TABLE users ADD COLUMN email VARCHAR(128)`); err != nil {
			return err
		}

		if err := m.ExecStmt(tx, "CREATE UNIQUE INDEX email_on_users ON users(email);"); err != nil {
			return err
		}

		return nil
	})
}
