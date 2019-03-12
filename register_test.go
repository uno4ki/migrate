package migrate

import (
	"database/sql"
	"fmt"
	"testing"
)

func migrate_stub(tx *sql.Tx) error { return nil }

func create_migrations() migrations_t {
	var m migrations_t

	// Create simple migrations array
	m.Register(3, "Migrate3", migrate_stub)
	m.Register(1, "Migrate1", migrate_stub)
	m.Register(2, "Migrate2", migrate_stub)

	return m
}

func TestRegister(t *testing.T) {
	m := create_migrations()

	t.Run("Array length == 3", func(t *testing.T) {
		if len(m) != 3 {
			t.Errorf("got '%d' want '3'", len(m))
		}
	})

	t.Run("Correct migrations order", func(t *testing.T) {
		for i, num := range []int{1, 2, 3} {
			if m[i].name != fmt.Sprintf("Migrate%d", num) {
				t.Errorf("got '%s' want 'Migrate%d'", m[0].name, num)
			}
		}
	})
}

func TestSelectMigrations(t *testing.T) {
	m := create_migrations()

	t.Run("Proper migrations selected", func(t *testing.T) {
		list := m.select_migrations([]string{"Migrate2"})

		if len(list) != 2 {
			t.Errorf("got '%d' want '2'\n", len(list))
		}

		for i, num := range []int{1, 3} {
			if list[i].name != fmt.Sprintf("Migrate%d", num) {
				t.Errorf("got '%s' want 'Migrate%d'", m[0].name, num)
			}
		}
	})
}
