package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func New() (conn *sql.DB, err error) {
	m, err := migrate.New(
		"file://internal/sqlite/migrations",
		"sqlite3://./app.db",
	)
	if err != nil {
		return nil, fmt.Errorf("creating migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if err != nil {
			return
		}
		if srcErr != nil {
			err = fmt.Errorf("closing migration source: %w", srcErr)
			return
		}
		if dbErr != nil {
			err = fmt.Errorf("closing migration database: %w", dbErr)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	log.Println("Migrations applied successfully")

	conn, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
