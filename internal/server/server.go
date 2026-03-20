package server

import (
	"log"
	"net/http"
	"time"

	"github.com/pawelgrzybek/go-notes/internal/store"
)

type wrapperWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrapperWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func NewServer(s store.Store) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlerList(s))
	mux.HandleFunc("GET /{noteID}", handlerGet(s))
	mux.HandleFunc("POST /", handlerCreate(s))
	mux.HandleFunc("DELETE /{noteID}", handlerDeleteOne(s))
	mux.HandleFunc("DELETE /", handlerDeleteAll(s))

	return Logging(mux)
}
