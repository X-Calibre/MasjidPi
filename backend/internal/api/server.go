package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/components"
	"github.com/X-Calibre/MasjidPi/backend/internal/listen"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/economic"
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
type masjidBoardEconomicProvider interface{ EconomicIndicators() *economic.Indicators }

type Server struct {
	httpServer                  *http.Server
	logger                      *slog.Logger
	playback                    *playback.Manager
	listen                      *listen.Controller
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
	installed                   components.Installed
}

type Config struct {
	Address                  string
	Frontend                 string
	CatalogueFile            string
	CatalogueDataRoot        string
	PreferencesPath          string
	Installed                components.Installed
	MasjidBoardHierarchyPath string
	MasjidBoardCataloguePath string
	MasjidBoardScopePath     string
}

type Dependencies struct {
	Logger                 *slog.Logger
	Playback               *playback.Manager
	Listen                 *listen.Controller
	Streams                *stream.Store
	Favourites             *storage.Favourites
	Preferences            *storage.Preferences
	AudioDeviceState       *storage.AudioDeviceState
	MasjidBoardService     masjidBoardStatusProvider
	MasjidBoardMaintenance masjidBoardMaintenance
}

func New(config Config, dependencies Dependencies) *Server {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(config.Frontend))
	preferences := dependencies.Preferences
	if preferences == nil {
		preferences = storage.NewPreferences(config.PreferencesPath)
	}
	server := &Server{
		logger:                   dependencies.Logger,
		playback:                 dependencies.Playback,
		listen:                   dependencies.Listen,
		streams:                  dependencies.Streams,
		favourites:               dependencies.Favourites,
		preferences:              preferences,
		audioDeviceState:         dependencies.AudioDeviceState,
		masjidBoardService:       dependencies.MasjidBoardService,
		masjidBoardMaintenance:   dependencies.MasjidBoardMaintenance,
		masjidBoardHierarchyPath: config.MasjidBoardHierarchyPath,
		masjidBoardCataloguePath: config.MasjidBoardCataloguePath,
		masjidBoardScopePath:     config.MasjidBoardScopePath,
		catalogueFile:            config.CatalogueFile,
		catalogueDataRoot:        config.CatalogueDataRoot,
		installed:                config.Installed,
		httpServer:               &http.Server{Addr: config.Address, Handler: mux},
	}
	server.SetMasjidBoardService(dependencies.MasjidBoardService)

	mux.HandleFunc("/api/components", server.components)
	mux.HandleFunc("/api/version", server.version)
	if config.Installed.Listen {
		mux.HandleFunc("/api/player/play", server.play)
		mux.HandleFunc("/api/player/stop", server.stop)
		mux.HandleFunc("/api/player/status", server.status)
		mux.HandleFunc("/api/player/volume", server.volume)
		mux.HandleFunc("/api/streams", server.streamsList)
		mux.HandleFunc("/api/favourites", server.favouritesHandler)
		mux.HandleFunc("/api/preferences", server.preferencesHandler)
		mux.HandleFunc("/api/catalogue/update", server.updateCatalogue)
		mux.HandleFunc("/api/listen/status", server.listenStatus)
		mux.HandleFunc("/api/listen/selection", server.listenSelection)
		mux.HandleFunc("/api/listen/power", server.listenPower)
		mux.HandleFunc("/api/listen/volume", server.listenVolume)
		mux.HandleFunc("/api/listen/radio-delay", server.listenRadioDelay)
		mux.HandleFunc("/api/listen/radio-schedule", server.listenRadioSchedule)
		mux.HandleFunc("/api/listen/radio-mode", server.listenRadioMode)
		mux.HandleFunc("/api/listen/start", server.listenStart)
		mux.HandleFunc("/api/listen/stop", server.listenStop)
	}
	if config.Installed.Board {
		mux.HandleFunc("/api/masjidboard/status", server.masjidBoardStatus)
		mux.HandleFunc("/api/masjidboard/boards/refresh", server.masjidBoardBoardsRefresh)
		mux.HandleFunc("/api/masjidboard/display", server.masjidBoardDisplay)
		mux.HandleFunc("/api/masjidboard/hierarchy", server.masjidBoardHierarchy)
		mux.HandleFunc("/api/masjidboard/hierarchy/refresh", server.masjidBoardHierarchyRefresh)
		mux.HandleFunc("/api/masjidboard/scope", server.masjidBoardScope)
		mux.HandleFunc("/api/masjidboard/catalogue", server.masjidBoardCatalogue)
		mux.HandleFunc("/api/masjidboard/catalogue/refresh", server.masjidBoardCatalogueRefresh)
		mux.HandleFunc("/api/masjidboard/selection", server.masjidBoardSelection)
		mux.HandleFunc("/api/masjidboard/layout", server.masjidBoardLayout)
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/", "/index.html":
			if !config.Installed.Listen && config.Installed.Board {
				http.Redirect(w, r, "/masjidboard-config.html", http.StatusTemporaryRedirect)
				return
			}
		case "/masjidboard-config.html", "/masjidboard.html":
			if !config.Installed.Board {
				http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
				return
			}
		}
		fileServer.ServeHTTP(w, r)
	})
	return server
}

func (s *Server) SetAudioDeviceState(state *storage.AudioDeviceState) { s.audioDeviceState = state }
func (s *Server) SetMasjidBoardService(service masjidBoardStatusProvider) {
	s.masjidBoardService = service
	if manager, ok := service.(masjidBoardSelectionManager); ok {
		s.masjidBoardSelectionManager = manager
	} else {
		s.masjidBoardSelectionManager = nil
	}
}
func (s *Server) SetMasjidBoardMaintenance(service masjidBoardMaintenance) {
	s.masjidBoardMaintenance = service
}
func (s *Server) SetMasjidBoardConfigurationPaths(hierarchyPath, scopePath, cataloguePath string) {
	s.masjidBoardHierarchyPath = hierarchyPath
	s.masjidBoardScopePath = scopePath
	s.masjidBoardCataloguePath = cataloguePath
}
func (s *Server) SetMasjidBoardCataloguePath(path string) { s.masjidBoardCataloguePath = path }
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
