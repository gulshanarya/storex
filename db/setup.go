package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

var DB *sql.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")
	port := os.Getenv("DB_PORT")

	var connStr strings.Builder

	_, err := fmt.Fprintf(&connStr, "postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)

	if err != nil {
		return fmt.Errorf("failed in concatenation of db connstr: %v", err)
	}

	db, err := sql.Open("postgres", connStr.String())

	if err != nil {
		return fmt.Errorf("error in opening db %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error in pinging db %w", err)
	}
	DB = db
	log.Println("successfully connected to db")
	return RunMigrations()
}

func RunMigrations() error {
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///home/gulshan-kumar/Documents/storex/db/migration",
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("Migration init error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Migration run error: %v", err)
	}

	log.Println("Database migrations applied successfully.")
	return nil
}
