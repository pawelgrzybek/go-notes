package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
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

	log.Println("server listening on :8080")

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNoteServiceServer(grpcServer, &server.Server{Store: s})
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
