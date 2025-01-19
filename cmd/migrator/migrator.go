package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Gonnekone/cryptoCrons/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	var migrationsPath, direction string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&direction, "direction", "up", "name of migrations table")
	flag.Parse()

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	cfg := config.MustLoad()

	m, err := migrate.New("file://"+migrationsPath, cfg.Storage.DSN()+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		err = m.Up()

	case "down":
		err = m.Down()
	}
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes to apply")

			return
		}
		log.Fatal(err)
	}

	fmt.Println("migrations applied successfuly")
}
