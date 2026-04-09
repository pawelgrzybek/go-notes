package server_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelgrzybek/go-notes/internal/db"
	"github.com/pawelgrzybek/go-notes/internal/notes"
	"github.com/pawelgrzybek/go-notes/internal/server"
)

func setupServer(t *testing.T) *server.Server {
	t.Helper()

	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			note TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	return server.NewServer(notes.NewService(db.New(conn)))
}
