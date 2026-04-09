package server_test

import (
	"context"
	"testing"
	"time"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc"
)

type mockWatchStream struct {
	grpc.ServerStream
	ctx  context.Context
	sent []*pb.WatchNotesResponse
}

func (m *mockWatchStream) Send(resp *pb.WatchNotesResponse) error {
	m.sent = append(m.sent, resp)
	return nil
}

func (m *mockWatchStream) Context() context.Context {
	return m.ctx
}

func TestWatchNotes(t *testing.T) {
	t.Run("sends initial snapshot", func(t *testing.T) {
		s := setupServer(t)

		s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("hello")})

		ctx, cancel := context.WithCancel(context.Background())
		stream := &mockWatchStream{ctx: ctx}

		done := make(chan error, 1)
		go func() {
			done <- s.WatchNotes(&pb.WatchNotesRequest{}, stream)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()
		<-done

		if len(stream.sent) < 1 {
			t.Fatal("expected at least one message")
		}
		if len(stream.sent[0].Notes) != 1 {
			t.Fatalf("expected 1 note in initial snapshot, got %d", len(stream.sent[0].Notes))
		}
		if stream.sent[0].Notes[0].GetNote() != "hello" {
			t.Errorf("expected 'hello', got %q", stream.sent[0].Notes[0].GetNote())
		}
	})

	t.Run("sends update on notification", func(t *testing.T) {
		s := setupServer(t)

		ctx, cancel := context.WithCancel(context.Background())
		stream := &mockWatchStream{ctx: ctx}

		done := make(chan error, 1)
		go func() {
			done <- s.WatchNotes(&pb.WatchNotesRequest{}, stream)
		}()

		time.Sleep(50 * time.Millisecond)

		s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("new note")})

		time.Sleep(50 * time.Millisecond)
		cancel()
		<-done

		if len(stream.sent) < 2 {
			t.Fatalf("expected at least 2 messages (initial + update), got %d", len(stream.sent))
		}
		if len(stream.sent[0].Notes) != 0 {
			t.Errorf("expected 0 notes in initial snapshot, got %d", len(stream.sent[0].Notes))
		}
		if len(stream.sent[1].Notes) != 1 {
			t.Errorf("expected 1 note after create, got %d", len(stream.sent[1].Notes))
		}
	})

	t.Run("stops on client disconnect", func(t *testing.T) {
		s := setupServer(t)

		ctx, cancel := context.WithCancel(context.Background())
		stream := &mockWatchStream{ctx: ctx}

		done := make(chan error, 1)
		go func() {
			done <- s.WatchNotes(&pb.WatchNotesRequest{}, stream)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("WatchNotes did not return after context cancellation")
		}
	})
}
