package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pawelgrzybek/go-notes/internal/server"
	"github.com/pawelgrzybek/go-notes/internal/store"
)

func run() error {
	db, err := sql.Open("sqlite3", "./app.db")

	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	fmt.Println("Successfully connected to SQLite database")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			note TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	s := store.NewStore(db)
	h := server.NewServer(s)

	svr := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	log.Println("server listening on :8080")
	return svr.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
