package server

import (
	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) WatchNotes(req *pb.WatchNotesRequest, stream pb.NoteService_WatchNotesServer) error {
	notes, err := s.Store.List()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}
	if err := stream.Send(&pb.WatchNotesResponse{Notes: notes}); err != nil {
		return err
	}

	ch := s.Notifier.Subscribe()
	defer s.Notifier.Unsubscribe(ch)

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case _, ok := <-ch:
			if !ok {
				return nil
			}
			notes, err := s.Store.List()
			if err != nil {
				return status.Errorf(codes.Internal, "failed to list notes: %v", err)
			}
			if err := stream.Send(&pb.WatchNotesResponse{Notes: notes}); err != nil {
				return err
			}
		}
	}
}
