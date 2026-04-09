package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/db"
	"github.com/pawelgrzybek/go-notes/internal/interceptors"
	"github.com/pawelgrzybek/go-notes/internal/notes"
	"github.com/pawelgrzybek/go-notes/internal/server"
	"github.com/pawelgrzybek/go-notes/internal/sqlite"
)

func run() error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	conn, err := sqlite.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	svc := notes.NewService(db.New(conn))
	svr := server.NewServer(svc)

	log.Println("server listening on :8080")

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.LoggingUnaryInterceptor(logger)),
		grpc.StreamInterceptor(interceptors.LoggingStreamInterceptor(logger)),
	)
	pb.RegisterNoteServiceServer(grpcServer, svr)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
