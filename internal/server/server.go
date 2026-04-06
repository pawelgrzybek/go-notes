package server

import (
	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/db"
	"github.com/pawelgrzybek/go-notes/internal/notifier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedNoteServiceServer
	Notifier *notifier.Notifier
	q        db.Querier
}

func NewServer(q db.Querier, n *notifier.Notifier) *Server {
	return &Server{
		Notifier: n,
		q:        q,
	}
}

func noteToProto(n db.Note) *pb.Note {
	return &pb.Note{
		Id:        new(int32(n.ID)),
		Note:      new(n.Note),
		CreatedAt: timestamppb.New(n.CreatedAt),
	}
}
