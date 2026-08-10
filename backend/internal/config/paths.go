package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	AppRoot         string
	DataRoot        string
	ConfigRoot      string
	Config           string
	Catalogue        string
	PlaybackState    string
	AudioDeviceState string
	VolumeState      string
	FavouritesState  string
	PreferencesState string
	Frontend         string
	Version          string
}

func NewPaths(base string) Paths {
	return Paths{
		AppRoot:          base,
		DataRoot:         filepath.Join(base, "backend", "data"),
		ConfigRoot:       filepath.Join(base, "backend", "configs"),
		Config:           filepath.Join(base, "backend", "configs", "default.yaml"),
		Catalogue:        filepath.Join(base, "backend", "data", "catalogue.json"),
		PlaybackState:    filepath.Join(base, "backend", "data", "playback.json"),
		AudioDeviceState: filepath.Join(base, "backend", "data", "audio_device.json"),
		VolumeState:      filepath.Join(base, "backend", "data", "volume.json"),
		FavouritesState:  filepath.Join(base, "backend", "data", "favourites.json"),
		PreferencesState: filepath.Join(base, "backend", "data", "preferences.json"),
		Frontend:         filepath.Join(base, "frontend"),
		Version:          filepath.Join(base, "version.json"),
	}
}

func RuntimePaths() (Paths, error) {
	if home := os.Getenv("MASJIDPI_HOME"); home != "" {
		return Paths{
			AppRoot:          home,
			DataRoot:         "/var/lib/masjidpi",
			ConfigRoot:       "/etc/masjidpi",
			Config:           "/etc/masjidpi/config.yaml",
			Catalogue:        "/var/lib/masjidpi/catalogue.json",
			PlaybackState:    "/var/lib/masjidpi/playback.json",
			AudioDeviceState: "/var/lib/masjidpi/audio_device.json",
			VolumeState:      "/var/lib/masjidpi/volume.json",
			FavouritesState:  "/var/lib/masjidpi/favourites.json",
			PreferencesState: "/var/lib/masjidpi/preferences.json",
			Frontend:         filepath.Join(home, "frontend"),
			Version:          filepath.Join(home, "version.json"),
		}, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}
	projectRoot := filepath.Dir(wd)
	return NewPaths(projectRoot), nil
}
