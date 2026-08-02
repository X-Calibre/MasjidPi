package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	player     *player.MPV
}

// New creates a new HTTP server.
func New(
	addr string,
	logger *slog.Logger,
	player *player.MPV,
) *Server {
	mux := http.NewServeMux()

	server := &Server{
		logger: logger,
		player: player,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}

	mux.HandleFunc("/", server.home)
	mux.HandleFunc("/api/player/play", server.play)
	mux.HandleFunc("/api/player/stop", server.stop)

	return server
}

func (s *Server) Start() error {
	s.logger.Info(
		"Starting HTTP server",
		"address", s.httpServer.Addr,
	)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(
		version.AppName +
			" is running\nVersion: " +
			version.Version +
			"\n",
	))
}
