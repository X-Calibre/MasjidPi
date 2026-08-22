package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeMasjidBoardStatusProvider struct {
	configured bool
	selection  selection.State
	results    []masjidboardruntime.Result
}

func (f fakeMasjidBoardStatusProvider) Configured() bool                     { return f.configured }
func (f fakeMasjidBoardStatusProvider) Selection() selection.State           { return f.selection }
func (f fakeMasjidBoardStatusProvider) Results() []masjidboardruntime.Result { return f.results }

func TestMasjidBoardStatusReturnsUnconfiguredState(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/status", nil)
	res := httptest.NewRecorder()

	s.masjidBoardStatus(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	var got masjidBoardStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Configured {
		t.Fatal("expected configured=false")
	}
	if len(got.Boards) != 0 {
		t.Fatalf("expected no boards, got %d", len(got.Boards))
	}
}

func TestMasjidBoardStatusPreservesSelectionOrderAndStaleFlag(t *testing.T) {
	attempt := time.Date(2026, 8, 18, 19, 0, 0, 0, time.UTC)
	success := attempt.Add(-time.Hour)
	first := selection.Board{
		CatalogueID:      "masjidboardlive:brits-jamia",
		Provider:         "masjidboardlive",
		ExternalID:       "brits-jamia",
		Name:             "Brits Jamia Masjid",
		TimeZoneOffsetMS: 7200000,
	}
	second := selection.Board{
		CatalogueID:      "masjidboardlive:brits-taqwa",
		Provider:         "masjidboardlive",
		ExternalID:       "brits-taqwa",
		Name:             "Masjid Taqwa",
		TimeZoneOffsetMS: 7200000,
	}
	cached := model.Board{Identity: model.BoardIdentity{ID: "brits-jamia", Name: "Brits Jamia Masjid", TimeZone: "GMT+02:00"}}

	s := &Server{masjidBoardService: fakeMasjidBoardStatusProvider{
		configured: true,
		selection:  selection.State{Boards: []selection.Board{first, second}},
		results: []masjidboardruntime.Result{
			{
				Selection:            first,
				Board:                &cached,
				Status:               masjidboardruntime.StatusStale,
				LastAttempt:          attempt,
				LastSuccessfulUpdate: success,
				UpdateError:          errors.New("upstream unavailable"),
			},
			{
				Selection:   second,
				Status:      masjidboardruntime.StatusUnavailable,
				LastAttempt: attempt,
				UpdateError: errors.New("no live data"),
			},
		},
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/status", nil)
	res := httptest.NewRecorder()
	s.masjidBoardStatus(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	var got masjidBoardStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !got.Configured || len(got.Boards) != 2 {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Boards[0].CatalogueID != first.CatalogueID || got.Boards[1].CatalogueID != second.CatalogueID {
		t.Fatalf("selection order changed: %+v", got.Boards)
	}
	if got.Boards[0].Status != masjidboardruntime.StatusStale || !got.Boards[0].UsingCachedData || !got.Boards[0].UpdateFailed {
		t.Fatalf("stale flags not exposed: %+v", got.Boards[0])
	}
	if got.Boards[0].Board == nil || got.Boards[0].UpdateError != "upstream unavailable" {
		t.Fatalf("stale board data/error missing: %+v", got.Boards[0])
	}
	if got.Boards[1].Status != masjidboardruntime.StatusUnavailable || got.Boards[1].UsingCachedData {
		t.Fatalf("unavailable state incorrect: %+v", got.Boards[1])
	}
}

func TestMasjidBoardStatusRejectsWrongMethod(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodPost, "/api/masjidboard/status", nil)
	res := httptest.NewRecorder()

	s.masjidBoardStatus(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}
