package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

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

	migration, err := os.ReadFile("internal/sql/migrations/schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migration))
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
