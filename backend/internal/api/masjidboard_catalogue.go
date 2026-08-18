package api

import (
	"net/http"
	"strings"
	"time"

	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

type masjidBoardCatalogueResponse struct {
	RetrievedAt *time.Time                       `json:"retrieved_at,omitempty"`
	ValidatedAt *time.Time                       `json:"validated_at,omitempty"`
	Count       int                              `json:"count"`
	Records     []masjidboardcatalogue.Record    `json:"records"`
}

func (s *Server) masjidBoardCatalogue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if strings.TrimSpace(s.masjidBoardCataloguePath) == "" {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard catalogue is not configured")
		return
	}

	statusFilter := strings.TrimSpace(r.URL.Query().Get("status"))
	if statusFilter != "" && !validMasjidBoardCatalogueStatus(statusFilter) {
		writeError(w, http.StatusBadRequest, "invalid MasjidBoard catalogue status")
		return
	}

	catalogue, err := masjidboardcatalogue.NewStore(s.masjidBoardCataloguePath).Load()
	if err != nil {
		s.logger.Warn("Could not load MasjidBoard catalogue", "error", err)
		writeError(w, http.StatusInternalServerError, "could not load MasjidBoard catalogue")
		return
	}

	query := r.URL.Query()
	q := normaliseFilter(query.Get("q"))
	name := normaliseFilter(query.Get("name"))
	city := normaliseFilter(query.Get("city"))
	region := normaliseFilter(query.Get("region"))
	if region == "" {
		region = normaliseFilter(query.Get("province"))
	}
	country := normaliseFilter(query.Get("country"))
	status := normaliseFilter(statusFilter)

	records := make([]masjidboardcatalogue.Record, 0, len(catalogue.Records))
	for _, record := range catalogue.Records {
		if !matchesMasjidBoardCatalogueRecord(record, q, name, city, region, country, status) {
			continue
		}
		records = append(records, record)
	}

	response := masjidBoardCatalogueResponse{
		Count:   len(records),
		Records: records,
	}
	if !catalogue.RetrievedAt.IsZero() {
		t := catalogue.RetrievedAt
		response.RetrievedAt = &t
	}
	if !catalogue.ValidatedAt.IsZero() {
		t := catalogue.ValidatedAt
		response.ValidatedAt = &t
	}

	writeJSON(w, http.StatusOK, response)
}

func matchesMasjidBoardCatalogueRecord(record masjidboardcatalogue.Record, q, name, city, region, country, status string) bool {
	if q != "" && !containsAny(q, record.Name, record.City, record.Region, record.Country) {
		return false
	}
	if name != "" && !containsFold(record.Name, name) {
		return false
	}
	if city != "" && !containsFold(record.City, city) {
		return false
	}
	if region != "" && !containsFold(record.Region, region) {
		return false
	}
	if country != "" && !containsFold(record.Country, country) {
		return false
	}
	if status != "" && normaliseFilter(string(record.Status)) != status {
		return false
	}
	return true
}

func validMasjidBoardCatalogueStatus(value string) bool {
	switch masjidboardcatalogue.Status(strings.ToLower(strings.TrimSpace(value))) {
	case masjidboardcatalogue.StatusActive, masjidboardcatalogue.StatusMissing, masjidboardcatalogue.StatusUnavailable:
		return true
	default:
		return false
	}
}

func containsAny(needle string, values ...string) bool {
	for _, value := range values {
		if containsFold(value, needle) {
			return true
		}
	}
	return false
}

func containsFold(value, needle string) bool {
	return strings.Contains(strings.ToLower(value), needle)
}

func normaliseFilter(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
