package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

type Server struct {
	httpServer  *http.Server
	logger      *slog.Logger
	playback    *playback.Manager
	streams     *stream.Store
	favourites  *storage.Favourites
	preferences *storage.Preferences
}

func New(
	addr string,
	logger *slog.Logger,
	playback *playback.Manager,
	streams *stream.Store,
	favourites *storage.Favourites,
	preferences *storage.Preferences,
	frontend string,
) *Server {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(frontend))

	server := &Server{
		logger:      logger,
		playback:    playback,
		streams:     streams,
		favourites:  favourites,
		preferences: preferences,
		httpServer: &http.Server{Addr: addr, Handler: mux},
	}

	mux.HandleFunc("/api/player/play", server.play)
	mux.HandleFunc("/api/player/stop", server.stop)
	mux.HandleFunc("/api/player/status", server.status)
	mux.HandleFunc("/api/player/volume", server.volume)
	mux.HandleFunc("/api/streams", server.streamsList)
	mux.HandleFunc("/api/favourites", server.favouritesHandler)
	mux.HandleFunc("/api/preferences", server.preferencesHandler)
	mux.HandleFunc("/api/catalogue/update", server.updateCatalogue)
	mux.HandleFunc("/api/version", server.version)
	mux.Handle("/", fileServer)

	return server
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", "address", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(version.AppName + " is running\nVersion: " + version.Version + "\n"))
}
