package server

import (
	"database/sql"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelgrzybek/go-notes/internal/store"
)

func setupTestServer(t *testing.T) (http.Handler, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			note TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	s := store.NewStore(db)
	h := NewServer(s)

	t.Cleanup(func() { db.Close() })

	return h, db
}

func insertNote(t *testing.T, db *sql.DB, note string) {
	t.Helper()
	_, err := db.Exec("INSERT INTO notes (note) VALUES (?)", note)
	if err != nil {
		t.Fatal(err)
	}
}
