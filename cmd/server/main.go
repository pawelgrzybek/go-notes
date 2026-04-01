package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/server"
	"github.com/pawelgrzybek/go-notes/internal/store"
)

func run() error {
	m, err := migrate.New(
		"file://internal/sql/migrations",
		"sqlite3://./app.db",
	)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return srcErr
	}
	if dbErr != nil {
		return dbErr
	}

	fmt.Println("Migrations applied successfully")

	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return err
	}
	defer db.Close()

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
