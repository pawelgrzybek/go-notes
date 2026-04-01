package store

import (
	"context"
	"database/sql"
	"errors"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/db"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store interface {
	List() ([]*pb.Note, error)
	Get(noteID int32) (*pb.Note, error)
	Create(note string) (*pb.Note, error)
	DeleteOne(noteID int32) (*pb.Note, error)
	DeleteAll() ([]*pb.Note, error)
}

type notesStore struct {
	q *db.Queries
}

func NewStore(conn *sql.DB) Store {
	return &notesStore{
		q: db.New(conn),
	}
}

func noteToProto(n db.Note) *pb.Note {
	id := int32(n.ID)
	note := n.Note
	return &pb.Note{
		Id:        &id,
		Note:      &note,
		CreatedAt: timestamppb.New(n.CreatedAt),
	}
}

func (s *notesStore) List() ([]*pb.Note, error) {
	rows, err := s.q.ListNotes(context.Background())
	if err != nil {
		return nil, err
	}

	notes := make([]*pb.Note, len(rows))
	for i, r := range rows {
		notes[i] = noteToProto(r)
	}

	return notes, nil
}

func (s *notesStore) Get(noteID int32) (*pb.Note, error) {
	row, err := s.q.GetNote(context.Background(), int64(noteID))
	if err != nil {
		return &pb.Note{}, err
	}

	return noteToProto(row), nil
}

func (s *notesStore) Create(n string) (*pb.Note, error) {
	row, err := s.q.CreateNote(context.Background(), n)
	if err != nil {
		return &pb.Note{}, err
	}

	return noteToProto(row), nil
}

func (s *notesStore) DeleteOne(noteID int32) (*pb.Note, error) {
	row, err := s.q.DeleteNoteByID(context.Background(), int64(noteID))
	if errors.Is(err, sql.ErrNoRows) {
		return &pb.Note{}, sql.ErrNoRows
	}
	if err != nil {
		return &pb.Note{}, err
	}

	return noteToProto(row), nil
}

func (s *notesStore) DeleteAll() ([]*pb.Note, error) {
	rows, err := s.q.DeleteAllNotes(context.Background())
	if err != nil {
		return nil, err
	}

	notes := make([]*pb.Note, len(rows))
	for i, r := range rows {
		notes[i] = noteToProto(r)
	}

	return notes, nil
}
