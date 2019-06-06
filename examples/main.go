// Sample
package main

import (
	"flag"
	_ "github.com/lib/pq"
	"github.com/uno4ki/migrate"
	"log"
)

var (
	dsn      string
	root_dsn string
	dbname   string
)

func main() {
	flag.StringVar(&dsn, "dsn", "", "Normal DSN (specify dbname here too)")
	flag.StringVar(&root_dsn, "root-dsn", "", "DSN for create/drop connection")
	flag.StringVar(&dbname, "dbname", "", "Sample database name")
	flag.Parse()

	if dsn == "" {
		log.Fatal("please specify -dsn")
	}

	if root_dsn == "" {
		log.Fatal("please specify -root-dsn")
	}

	if dbname == "" {
		log.Fatal("please specify -dbname")
	}

	m := migrate.NewMigrator("postgres", dsn, root_dsn, dbname)

	log.Print("==> create database")

	if err := m.Create(); err != nil {
		log.Fatalf("Failed to create database: %v\n", err)
	}

	log.Print("==> run database migrations")
	if err := m.Migrate(); err != nil {
		log.Fatalf("Failed to seed database: %v\n", err)
	}

	log.Print("==> seed database")
	if err := m.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v\n", err)
	}

	log.Print("==> drop database :)")
	if err := m.Drop(); err != nil {
		log.Fatalf("Failed to drop database: %v\n", err)
	}
}
