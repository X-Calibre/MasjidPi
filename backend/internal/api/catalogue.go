package api

import "github.com/X-Calibre/MasjidPi/backend/internal/catalogue"

func (s *Server) refreshCatalogue() error {
	streams, err := catalogue.Update(s.catalogueFile)
	if err != nil {
		return err
	}

	s.streams.Replace(streams)
	return nil
}
