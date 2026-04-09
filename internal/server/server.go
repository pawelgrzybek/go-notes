package server

import (
	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/notes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedNoteServiceServer
	svc notes.Service
}

func NewServer(svc notes.Service) *Server {
	return &Server{svc: svc}
}

func noteToProto(n notes.Note) *pb.Note {
	return &pb.Note{
		Id:        new(int32(n.ID)),
		Note:      new(n.Note),
		CreatedAt: timestamppb.New(n.CreatedAt),
	}
}
