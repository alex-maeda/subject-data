package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// Server wraps the HTTP server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// NewServer creates a new Server with the given router and port.
func NewServer(handler http.Handler, port string, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + port,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger,
	}
}

// Start begins serving HTTP requests. Blocks until the server stops.
func (s *Server) Start() error {
	s.logger.Info("subject-data listening", "port", s.httpServer.Addr[1:])
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down")
	return s.httpServer.Shutdown(ctx)
}
