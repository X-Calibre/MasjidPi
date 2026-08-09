package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	AppRoot    string
	DataRoot   string
	ConfigRoot string

	Config        string
	Catalogue     string
	PlaybackState string
	Frontend      string
	Version       string
}

func NewPaths(base string) Paths {
	return Paths{
		AppRoot:    base,
		DataRoot:   filepath.Join(base, "backend", "data"),
		ConfigRoot: filepath.Join(base, "backend", "configs"),

		Config:        filepath.Join(base, "backend", "configs", "default.yaml"),
		Catalogue:     filepath.Join(base, "backend", "data", "catalogue.json"),
		PlaybackState: filepath.Join(base, "backend", "data", "playback.json"),
		Frontend:      filepath.Join(base, "frontend"),
		Version:       filepath.Join(base, "version.json"),
	}
}

func RuntimePaths() (Paths, error) {

	// Installed runtime
	if home := os.Getenv("MASJIDPI_HOME"); home != "" {

		return Paths{
			AppRoot:    home,
			DataRoot:   "/var/lib/masjidpi",
			ConfigRoot: "/etc/masjidpi",

			Config:        "/etc/masjidpi/config.yaml",
			Catalogue:     "/var/lib/masjidpi/catalogue.json",
			PlaybackState: "/var/lib/masjidpi/playback.json",
			Frontend:      filepath.Join(home, "frontend"),
			Version:       filepath.Join(home, "version.json"),
		}, nil
	}

	// Development runtime
	wd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}

	projectRoot := filepath.Dir(wd)

	paths := Paths{
		AppRoot:    projectRoot,
		DataRoot:   filepath.Join(projectRoot, "backend", "data"),
		ConfigRoot: filepath.Join(projectRoot, "backend", "configs"),

		Config:        filepath.Join(projectRoot, "backend", "configs", "default.yaml"),
		Catalogue:     filepath.Join(projectRoot, "backend", "data", "catalogue.json"),
		PlaybackState: filepath.Join(projectRoot, "backend", "data", "playback.json"),
		Frontend:      filepath.Join(projectRoot, "frontend"),
		Version:       filepath.Join(projectRoot, "version.json"),
	}

	return paths, nil
}
