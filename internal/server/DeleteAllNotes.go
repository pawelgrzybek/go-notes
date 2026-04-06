package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func (s *Server) DeleteAllNotes(ctx context.Context, req *pb.DeleteAllNotesRequest) (*pb.DeleteAllNotesResponse, error) {
	rows, err := s.q.DeleteAllNotes(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete all notes: %v", err)
	}

	notes := make([]*pb.Note, len(rows))
	for i, r := range rows {
		notes[i] = noteToProto(r)
	}

	if s.Notifier != nil {
		s.Notifier.Notify()
	}

	return &pb.DeleteAllNotesResponse{
		Notes: notes,
	}, nil
}
