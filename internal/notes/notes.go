package notes

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/pawelgrzybek/go-notes/internal/db"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        int64
	Note      string
	CreatedAt time.Time
}

type Service interface {
	Create(ctx context.Context, content string) (Note, error)
	Get(ctx context.Context, id int64) (Note, error)
	List(ctx context.Context) ([]Note, error)
	Delete(ctx context.Context, id int64) (Note, error)
	DeleteAll(ctx context.Context) ([]Note, error)
	Watch(ctx context.Context) <-chan struct{}
}

type service struct {
	q           db.Querier
	mu          sync.Mutex
	subscribers map[chan struct{}]struct{}
}

func NewService(q db.Querier) Service {
	return &service{
		q:           q,
		subscribers: make(map[chan struct{}]struct{}),
	}
}

func fromDB(n db.Note) Note {
	return Note{
		ID:        n.ID,
		Note:      n.Note,
		CreatedAt: n.CreatedAt,
	}
}

func fromDBAll(rows []db.Note) []Note {
	out := make([]Note, len(rows))
	for i, r := range rows {
		out[i] = fromDB(r)
	}
	return out
}

func (s *service) Create(ctx context.Context, content string) (Note, error) {
	row, err := s.q.CreateNote(ctx, content)
	if err != nil {
		return Note{}, err
	}
	s.notify()
	return fromDB(row), nil
}

func (s *service) Get(ctx context.Context, id int64) (Note, error) {
	row, err := s.q.GetNote(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return fromDB(row), nil
}

func (s *service) List(ctx context.Context) ([]Note, error) {
	rows, err := s.q.ListNotes(ctx)
	if err != nil {
		return nil, err
	}
	return fromDBAll(rows), nil
}

func (s *service) Delete(ctx context.Context, id int64) (Note, error) {
	row, err := s.q.DeleteNoteByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	s.notify()
	return fromDB(row), nil
}

func (s *service) DeleteAll(ctx context.Context) ([]Note, error) {
	rows, err := s.q.DeleteAllNotes(ctx)
	if err != nil {
		return nil, err
	}
	s.notify()
	return fromDBAll(rows), nil
}

// Watch returns a signal channel that receives a value each time the notes
// collection changes. Subscribers should re-fetch the list with List on every
// signal. The channel is closed when ctx is cancelled.
//
// The channel has a buffer of 1 and notifications are delivered with a
// non-blocking send: if a signal is already queued it is coalesced, since a
// single wake-up is enough for the subscriber to observe the latest state.
func (s *service) Watch(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{}, 1)

	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.subscribers, ch)
		close(ch)
		s.mu.Unlock()
	}()

	return ch
}

func (s *service) notify() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
