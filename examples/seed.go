// Sample SEED file
package main

import (
	"database/sql"
	"fmt"
	m "github.com/uno4ki/migrate"
)

func init() {

	// Create users in one call
	m.Seeds.Register(1, "create users", func(tx *sql.Tx) error {
		return m.ExecStmt(tx, `INSERT INTO users (name, email)
			VALUES('name1', 'mail1@example.com'),
					  ('name2', 'mail2@example.com'),
						('name3', 'mail3@example.com');`)
	})

	// Create multiple users from array
	m.Seeds.Register(2, "multiple users", func(tx *sql.Tx) error {
		for _, v := range []string{"example1@example.com", "example2@example.com", "example3@example.com"} {
			if err := m.ExecStmt(tx, fmt.Sprintf(`INSERT INTO users (name, email) VALUES('say my name', '%s');`, v)); err != nil {
				return err
			}
		}
		return nil
	})
}
