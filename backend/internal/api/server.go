package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

type masjidBoardStatusProvider interface {
	Configured() bool
	Selection() selection.State
	Results() []masjidboardruntime.Result
}

type Server struct {
	httpServer                  *http.Server
	logger                      *slog.Logger
	playback                    *playback.Manager
	streams                     *stream.Store
	favourites                  *storage.Favourites
	preferences                 *storage.Preferences
	audioDeviceState            *storage.AudioDeviceState
	masjidBoardService          masjidBoardStatusProvider
	masjidBoardSelectionManager masjidBoardSelectionManager
	masjidBoardMaintenance      masjidBoardMaintenance
	masjidBoardHierarchyPath    string
	masjidBoardCataloguePath    string
	masjidBoardScopePath        string
	catalogueFile               string
	catalogueDataRoot           string
}

func New(addr string, logger *slog.Logger, playback *playback.Manager, streams *stream.Store, favourites *storage.Favourites, frontend, catalogueFile, catalogueDataRoot string) *Server {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(frontend))

	preferencesPath := "/var/lib/masjidpi/preferences.json"
	if home := os.Getenv("MASJIDPI_HOME"); home != "" {
		preferencesPath = filepath.Join(home, "backend", "data", "preferences.json")
	}

	server := &Server{
		logger:            logger,
		playback:          playback,
		streams:           streams,
		favourites:        favourites,
		preferences:       storage.NewPreferences(preferencesPath),
		catalogueFile:     catalogueFile,
		catalogueDataRoot: catalogueDataRoot,
		httpServer:        &http.Server{Addr: addr, Handler: mux},
	}

	mux.HandleFunc("/api/components", server.components)
	mux.HandleFunc("/api/player/play", server.play)
	mux.HandleFunc("/api/player/stop", server.stop)
	mux.HandleFunc("/api/player/status", server.status)
	mux.HandleFunc("/api/player/volume", server.volume)
	mux.HandleFunc("/api/streams", server.streamsList)
	mux.HandleFunc("/api/favourites", server.favouritesHandler)
	mux.HandleFunc("/api/preferences", server.preferencesHandler)
	mux.HandleFunc("/api/catalogue/update", server.updateCatalogue)
	mux.HandleFunc("/api/masjidboard/status", server.masjidBoardStatus)
	mux.HandleFunc("/api/masjidboard/boards/refresh", server.masjidBoardBoardsRefresh)
	mux.HandleFunc("/api/masjidboard/display", server.masjidBoardDisplay)
	mux.HandleFunc("/api/masjidboard/hierarchy", server.masjidBoardHierarchy)
	mux.HandleFunc("/api/masjidboard/hierarchy/refresh", server.masjidBoardHierarchyRefresh)
	mux.HandleFunc("/api/masjidboard/scope", server.masjidBoardScope)
	mux.HandleFunc("/api/masjidboard/catalogue", server.masjidBoardCatalogue)
	mux.HandleFunc("/api/masjidboard/catalogue/refresh", server.masjidBoardCatalogueRefresh)
	mux.HandleFunc("/api/masjidboard/selection", server.masjidBoardSelection)
	mux.HandleFunc("/api/version", server.version)
	mux.Handle("/", fileServer)

	return server
}

func (s *Server) SetAudioDeviceState(state *storage.AudioDeviceState) {
	s.audioDeviceState = state
}

// SetMasjidBoardService retains the MasjidBoard runtime status provider without
// coupling MasjidBoard availability to the audio subsystem. Production service
// implementations may also support live selection reconfiguration.
func (s *Server) SetMasjidBoardService(service masjidBoardStatusProvider) {
	s.masjidBoardService = service
	if manager, ok := service.(masjidBoardSelectionManager); ok {
		s.masjidBoardSelectionManager = manager
	} else {
		s.masjidBoardSelectionManager = nil
	}
}

// SetMasjidBoardMaintenance exposes explicit hierarchy/catalogue maintenance
// operations to the configuration API.
func (s *Server) SetMasjidBoardMaintenance(service masjidBoardMaintenance) {
	s.masjidBoardMaintenance = service
}

// SetMasjidBoardConfigurationPaths configures the disk-first discovery state
// read and written by the WebUI/API configuration surface.
func (s *Server) SetMasjidBoardConfigurationPaths(hierarchyPath, scopePath, cataloguePath string) {
	s.masjidBoardHierarchyPath = hierarchyPath
	s.masjidBoardScopePath = scopePath
	s.masjidBoardCataloguePath = cataloguePath
}

// SetMasjidBoardCataloguePath is retained for callers/tests that only need the
// catalogue read API.
func (s *Server) SetMasjidBoardCataloguePath(path string) {
	s.masjidBoardCataloguePath = path
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
