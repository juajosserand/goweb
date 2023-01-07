package httpserver

import (
	"context"
	"net/http"
	"time"
)

type server struct {
	server  *http.Server
	timeout time.Duration
	notify  chan error
}

func New(h http.Handler, ops ...Option) *server {
	httpServer := &http.Server{
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         ":8080",
	}

	s := &server{
		server:  httpServer,
		timeout: 3 * time.Second,
		notify:  make(chan error),
	}

	for _, op := range ops {
		op(s)
	}

	s.run()

	return s
}

func (s *server) run() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *server) Notify() <-chan error {
	return s.notify
}

func (s *server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
