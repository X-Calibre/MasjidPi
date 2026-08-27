package api

import (
	"github.com/X-Calibre/MasjidPi/backend/internal/catalogue"
	"github.com/X-Calibre/MasjidPi/backend/internal/radio"
)

func (s *Server) refreshCatalogue() error {
	streams, err := catalogue.Update(s.catalogueFile)
	if err != nil {
		return err
	}

	s.streams.Replace(radio.Merge(streams))
	return nil
}
