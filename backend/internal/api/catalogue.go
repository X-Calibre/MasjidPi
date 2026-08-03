package api

import "github.com/X-Calibre/MasjidPi/backend/internal/catalogue"

func (s *Server) refreshCatalogue() error {

	if err := catalogue.Update(); err != nil {
		return err
	}

	return s.streams.Reload(catalogueFile)
}
