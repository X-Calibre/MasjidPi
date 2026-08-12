package api

import (
	"path/filepath"

	"github.com/X-Calibre/MasjidPi/backend/internal/catalogue"
)

func (s *Server) refreshCatalogue() error {
	if err := catalogue.Update(filepath.Join(s.catalogueDataRoot, "page.html"), s.catalogueFile); err != nil {
		return err
	}

	return s.streams.Reload(s.catalogueFile)
}
