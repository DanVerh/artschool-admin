package application

import (
	"log"
	"fmt"
	"os"
	"reflect"
	"strings"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
)

// Create application struct (class) with required for migration fields
type App struct {
	File string
	DbUri string
}

// Construct for the App object
func New() *App {
	app := &App{}

	// Set App fields with env vars
	val := reflect.ValueOf(app).Elem()
	for i := 0; i < val.NumField(); i++ {
		// Get the field name in uppercase
		field := strings.ToUpper(val.Type().Field(i).Name)
		// file field is hardcoded and not set with env var
		if field == "FILE" {
			app.File = "file://migrations"
			continue
		}
		// Check if env var is set
		if os.Getenv(field) == "" {
			log.Fatal(field, " env var is missing")
		}
		// Set the obect field value
		val.Field(i).SetString(os.Getenv(field))
	}

	return app
}

// Start the app by running the migration
func (app *App) Start() error {
	m, err := migrate.New(app.File, app.DbUri)
    if err != nil {
        log.Fatalf("Failed to create migrate instance: %v", err)
    }

	defer func () {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Fatalf("Migration source closure error: %v", srcErr)
		}
		if dbErr != nil {
			log.Fatalf("DB connection closure error: %v", dbErr)
		}
	} ()

	err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Failed to run up migrations: %v", err)
    }

	version, _, err := m.Version()
	if err != nil {
    	fmt.Errorf("Failed to get current migration version: %v", err)
	}
	fmt.Printf("Migrations applied successfully. Current version is %v\n", version)

	return nil
}