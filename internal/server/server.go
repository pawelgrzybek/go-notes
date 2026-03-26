package server

import (
	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/store"
)

type Server struct {
	pb.UnimplementedNoteServiceServer
	Store store.Store
}
