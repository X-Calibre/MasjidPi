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

	// Installed runtime
	if home := os.Getenv("MASJIDPI_HOME"); home != "" {
		return NewPaths(home), nil
	}

	// Development runtime
	wd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}

	projectRoot := filepath.Dir(wd)

	paths := NewPaths(projectRoot)

	// Development layout
	paths.Config = filepath.Join(projectRoot, "backend", "configs", "default.yaml")
	paths.Catalogue = filepath.Join(projectRoot, "backend", "data", "catalogue.json")
	paths.Frontend = filepath.Join(projectRoot, "frontend")

	return paths, nil
}
