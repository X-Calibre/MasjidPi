package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Base      string
	Config    string
	Catalogue string
	Frontend  string
}

func NewPaths(base string) Paths {
	return Paths{
		Base:      base,
		Config:    filepath.Join(base, "configs", "default.yaml"),
		Catalogue: filepath.Join(base, "data", "catalogue.json"),
		Frontend:  filepath.Join(base, "frontend"),
	}
}

func RuntimePaths() (Paths, error) {

	// If MasjidPi is installed, the installer or systemd service
	// will set MASJIDPI_HOME.
	if home := os.Getenv("MASJIDPI_HOME"); home != "" {
		return NewPaths(home), nil
	}

	// Development mode (running from backend/)
	wd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}

	return NewPaths(wd), nil
}
