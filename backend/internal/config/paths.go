package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	AppRoot              string
	DataRoot             string
	ConfigRoot           string
	Config               string
	Catalogue            string
	PlaybackState        string
	AudioDeviceState     string
	VolumeState          string
	FavouritesState      string
	PreferencesState     string
	DisplaySettingsState string
	TimezoneState        string
	UpdateState          string
	UpdateDownloads      string
	MasjidBoardHierarchy string
	MasjidBoardCatalogue string
	MasjidBoardScope     string
	MasjidBoardSelection string
	MasjidBoardCache     string
	Frontend             string
	Version              string
}

func NewPaths(base string) Paths {
	return Paths{
		AppRoot:              base,
		DataRoot:             filepath.Join(base, "backend", "data"),
		ConfigRoot:           filepath.Join(base, "backend", "configs"),
		Config:               filepath.Join(base, "backend", "configs", "default.yaml"),
		Catalogue:            filepath.Join(base, "backend", "data", "catalogue.json"),
		PlaybackState:        filepath.Join(base, "backend", "data", "playback.json"),
		AudioDeviceState:     filepath.Join(base, "backend", "data", "audio_device.json"),
		VolumeState:          filepath.Join(base, "backend", "data", "volume.json"),
		FavouritesState:      filepath.Join(base, "backend", "data", "favourites.json"),
		PreferencesState:     filepath.Join(base, "backend", "data", "preferences.json"),
		DisplaySettingsState: filepath.Join(base, "backend", "data", "display_settings.json"),
		TimezoneState:        filepath.Join(base, "backend", "data", "timezone.json"),
		UpdateState:          filepath.Join(base, "backend", "data", "update_state.json"),
		UpdateDownloads:      filepath.Join(base, "backend", "data", "update-downloads"),
		MasjidBoardHierarchy: filepath.Join(base, "backend", "data", "masjidboard_hierarchy.json"),
		MasjidBoardCatalogue: filepath.Join(base, "backend", "data", "masjidboard_catalogue.json"),
		MasjidBoardScope:     filepath.Join(base, "backend", "data", "masjidboard_scope.json"),
		MasjidBoardSelection: filepath.Join(base, "backend", "data", "masjidboard_selection.json"),
		MasjidBoardCache:     filepath.Join(base, "backend", "data", "masjidboard_cache"),
		Frontend:             filepath.Join(base, "frontend"),
		Version:              filepath.Join(base, "version.json"),
	}
}

func RuntimePaths() (Paths, error) {
	if home := os.Getenv("MASJIDFRAME_HOME"); home != "" {
		return Paths{
			AppRoot:              home,
			DataRoot:             "/var/lib/masjidframe",
			ConfigRoot:           "/etc/masjidframe",
			Config:               "/etc/masjidframe/config.yaml",
			Catalogue:            "/var/lib/masjidframe/catalogue.json",
			PlaybackState:        "/var/lib/masjidframe/playback.json",
			AudioDeviceState:     "/var/lib/masjidframe/audio_device.json",
			VolumeState:          "/var/lib/masjidframe/volume.json",
			FavouritesState:      "/var/lib/masjidframe/favourites.json",
			PreferencesState:     "/var/lib/masjidframe/preferences.json",
			DisplaySettingsState: "/var/lib/masjidframe/display_settings.json",
			TimezoneState:        "/var/lib/masjidframe/timezone.json",
			UpdateState:          "/var/lib/masjidframe/update_state.json",
			UpdateDownloads:      "/var/lib/masjidframe/update-downloads",
			MasjidBoardHierarchy: "/var/lib/masjidframe/masjidboard_hierarchy.json",
			MasjidBoardCatalogue: "/var/lib/masjidframe/masjidboard_catalogue.json",
			MasjidBoardScope:     "/var/lib/masjidframe/masjidboard_scope.json",
			MasjidBoardSelection: "/var/lib/masjidframe/masjidboard_selection.json",
			MasjidBoardCache:     "/var/lib/masjidframe/masjidboard_cache",
			Frontend:             filepath.Join(home, "frontend"),
			Version:              filepath.Join(home, "version.json"),
		}, nil
	}

	if executable, err := os.Executable(); err == nil {
		executableDir := filepath.Dir(executable)
		packagedConfig := filepath.Join(executableDir, "default.yaml")
		if _, err := os.Stat(packagedConfig); err == nil {
			return Paths{
				AppRoot:              executableDir,
				DataRoot:             "/var/lib/masjidframe",
				ConfigRoot:           executableDir,
				Config:               packagedConfig,
				Catalogue:            "/var/lib/masjidframe/catalogue.json",
				PlaybackState:        "/var/lib/masjidframe/playback.json",
				AudioDeviceState:     "/var/lib/masjidframe/audio_device.json",
				VolumeState:          "/var/lib/masjidframe/volume.json",
				FavouritesState:      "/var/lib/masjidframe/favourites.json",
				PreferencesState:     "/var/lib/masjidframe/preferences.json",
				DisplaySettingsState: "/var/lib/masjidframe/display_settings.json",
				TimezoneState:        "/var/lib/masjidframe/timezone.json",
				UpdateState:          "/var/lib/masjidframe/update_state.json",
				UpdateDownloads:      "/var/lib/masjidframe/update-downloads",
				MasjidBoardHierarchy: "/var/lib/masjidframe/masjidboard_hierarchy.json",
				MasjidBoardCatalogue: "/var/lib/masjidframe/masjidboard_catalogue.json",
				MasjidBoardScope:     "/var/lib/masjidframe/masjidboard_scope.json",
				MasjidBoardSelection: "/var/lib/masjidframe/masjidboard_selection.json",
				MasjidBoardCache:     "/var/lib/masjidframe/masjidboard_cache",
				Frontend:             filepath.Join(executableDir, "frontend"),
				Version:              filepath.Join(executableDir, "VERSION"),
			}, nil
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}
	projectRoot := filepath.Dir(wd)
	return NewPaths(projectRoot), nil
}
