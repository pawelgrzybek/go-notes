package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func (s *Server) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	result, err := s.svc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	protoNotes := make([]*pb.Note, len(result))
	for i, n := range result {
		protoNotes[i] = noteToProto(n)
	}

	return &pb.ListNotesResponse{
		Notes: protoNotes,
	}, nil
}
