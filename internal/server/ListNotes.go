package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func (s *Server) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	rows, err := s.q.ListNotes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	notes := make([]*pb.Note, len(rows))
	for i, r := range rows {
		notes[i] = noteToProto(r)
	}

	return &pb.ListNotesResponse{
		Notes: notes,
	}, nil

}
