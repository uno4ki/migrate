package migrate

import (
	"database/sql"
	"log"
	"sort"
)

// callback function types
type migrate_t func(*sql.Tx) error

/**
 * Migration struct desription.
 */
type migration_t struct {
	order   int
	name    string
	execute migrate_t
}

// List of migration structs
type migrations_t []migration_t

// All registered migrations
var Migrations migrations_t
var Seeds migrations_t

// "sort" stuff
func (v migrations_t) Len() int {
	return len(v)
}

func (v migrations_t) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v migrations_t) Less(i, j int) bool {
	if v[i].order == v[j].order {
		log.Fatalf("Same order migrations found: [%s:%d] == [%s:%d]",
			v[i].name, v[i].order, v[j].name, v[j].order)
	}

	return v[i].order < v[j].order
}

// Each migration should be registered with these function.
// All migrations sorted by the order field
//
// @param fun   - migration function ptr
// @param name  - migration name
// @param order - unique migration order (same orders are error)
func (list *migrations_t) Register(order int, name string, execute migrate_t) {
	m := migration_t{order, name, execute}
	*list = append(*list, m)
	sort.Sort(list)
}

func (list *migrations_t) select_migrations(applied []string) migrations_t {
	var ret migrations_t

	for _, m := range *list {
		found := false
		for _, a := range applied {
			if a == m.name {
				found = true
				break
			}
		}

		if found {
			log.Printf("--- Skip migration: %s (already migrated)\n", m.name)
		} else {
			log.Printf("+++ Queue migration: %s\n", m.name)
			ret = append(ret, m)
		}
	}

	sort.Sort(ret)

	return ret
}

// Execute migration statement helper
func ExecStmt(tx *sql.Tx, statement string) error {

	if stmt, err := tx.Prepare(statement); err != nil {
		return err
	} else {
		if _, err := stmt.Exec(); err != nil {
			return err
		}
		stmt.Close()
	}

	return nil
}
