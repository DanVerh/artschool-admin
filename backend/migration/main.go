package main

import (
	"log"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
)

func main() {
	m, err := migrate.New(
        "file://migrations",
        "mongodb://localhost:27017/artschool-admin",
    )
    if err != nil {
        log.Fatalf("Failed to create migrate instance: %v", err)
    }

	err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Failed to run up migrations: %v", err)
    }
    fmt.Println("Migrations applied successfully")
}