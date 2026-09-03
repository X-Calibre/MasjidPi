package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/display"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func TestMasjidBoardDisplayReturnsPresentationOnlyView(t *testing.T) {
	first := selection.Board{CatalogueID: "masjidboardlive:one", Provider: "masjidboardlive", ExternalID: "one", Name: "One", TimeZoneOffsetMS: 7200000}
	second := selection.Board{CatalogueID: "masjidboardlive:two", Provider: "masjidboardlive", ExternalID: "two", Name: "Two", TimeZoneOffsetMS: 7200000}
	updated := time.Date(2026, 8, 19, 19, 0, 0, 0, time.UTC)
	cached := model.Board{
		Identity:      model.BoardIdentity{ID: "one", Name: "One Masjid", TimeZone: "GMT+02:00"},
		PrayerTimes:   model.PrayerTimes{Asr: model.PrayerTime{Jamaah: &model.ClockTime{Hour: 16, Minute: 45}}},
		Announcements: []model.Announcement{{Title: "Masjid announcement", Content: "Programme after Esha"}},
		Notices: []model.Notice{{
			Type: model.NoticeTypeNikah, Title: "Nikah Notice", Fields: map[string]string{"date": "2026-08-29"},
		}},
	}

	s := &Server{masjidBoardService: fakeMasjidBoardStatusProvider{
		configured: true,
		selection:  selection.State{Theme: selection.ThemeRuby, SlideDurationSeconds: 30, Boards: []selection.Board{first, second}},
		results: []masjidboardruntime.Result{
			{Selection: first, Board: &cached, Status: masjidboardruntime.StatusStale, LastSuccessfulUpdate: updated, UpdateError: errors.New("secret diagnostic")},
			{Selection: second, Status: masjidboardruntime.StatusUnavailable, UpdateError: errors.New("another diagnostic")},
		},
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/display", nil)
	res := httptest.NewRecorder()
	s.masjidBoardDisplay(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	body := res.Body.String()
	var got display.View
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !got.Configured || len(got.Boards) != 2 {
		t.Fatalf("view = %+v", got)
	}
	if got.Theme != selection.ThemeRuby || got.SlideDuration != 30 {
		t.Fatalf("display preferences = %+v", got)
	}
	if strings.Contains(body, `"layout"`) {
		t.Fatalf("display API still exposes removed layout preference: %s", body)
	}
	if got.Boards[0].CatalogueID != first.CatalogueID || got.Boards[1].CatalogueID != second.CatalogueID {
		t.Fatalf("order changed: %+v", got.Boards)
	}
	if !got.Boards[0].Stale || got.Boards[0].Status != masjidboardruntime.StatusStale {
		t.Fatalf("stale state = %+v", got.Boards[0])
	}
	if len(got.Boards[0].Prayers) != 5 || got.Boards[0].Prayers[2].Jamaah == nil || got.Boards[0].Prayers[2].Jamaah.Minute != 45 {
		t.Fatalf("prayer presentation = %+v", got.Boards[0].Prayers)
	}
	if got.Boards[1].Status != masjidboardruntime.StatusUnavailable {
		t.Fatalf("unavailable board = %+v", got.Boards[1])
	}
	if len(got.Boards[0].Announcements) != 1 || got.Boards[0].Announcements[0].Title != "Masjid announcement" {
		t.Fatalf("announcement presentation = %+v", got.Boards[0].Announcements)
	}
	if len(got.Boards[0].Notices) != 1 || got.Boards[0].Notices[0].Type != string(model.NoticeTypeNikah) {
		t.Fatalf("notice presentation = %+v", got.Boards[0].Notices)
	}
	if strings.Contains(body, "secret diagnostic") || strings.Contains(body, "another diagnostic") || strings.Contains(body, "\"provider\"") || strings.Contains(body, "external_id") {
		t.Fatalf("display API leaked administrative/provider detail: %s", body)
	}
	if strings.Contains(body, "last_successful_update") {
		t.Fatalf("display API leaked refresh metadata that destabilises its ETag: %s", body)
	}
}

func TestMasjidBoardDisplayETagIgnoresSuccessfulRefreshTimestamp(t *testing.T) {
	selectedBoard := selection.Board{CatalogueID: "masjidboardlive:one", Name: "One"}
	board := model.Board{Identity: model.BoardIdentity{ID: "one", Name: "One Masjid"}}
	selectionState := selection.State{Boards: []selection.Board{selectedBoard}}
	response := func(updated time.Time) *httptest.ResponseRecorder {
		s := &Server{masjidBoardService: fakeMasjidBoardStatusProvider{
			configured: true,
			selection:  selectionState,
			results: []masjidboardruntime.Result{{
				Selection: selectedBoard, Board: &board, Status: masjidboardruntime.StatusCurrent, LastSuccessfulUpdate: updated,
			}},
		}}
		res := httptest.NewRecorder()
		s.masjidBoardDisplay(res, httptest.NewRequest(http.MethodGet, "/api/masjidboard/display", nil))
		return res
	}

	first := response(time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC))
	second := response(time.Date(2026, 9, 3, 8, 30, 0, 0, time.UTC))
	if first.Header().Get("ETag") != second.Header().Get("ETag") {
		t.Fatalf("unchanged display content received different ETags: %q and %q", first.Header().Get("ETag"), second.Header().Get("ETag"))
	}
	if first.Body.String() != second.Body.String() {
		t.Fatal("unchanged display content changed when only refresh metadata advanced")
	}
}

func TestMasjidBoardDisplayReturnsUnconfiguredState(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/display", nil)
	res := httptest.NewRecorder()
	s.masjidBoardDisplay(res, req)

	var got display.View
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Configured || len(got.Boards) != 0 {
		t.Fatalf("view = %+v", got)
	}
}

func TestMasjidBoardDisplaySupportsConditionalRequests(t *testing.T) {
	s := &Server{}
	firstRequest := httptest.NewRequest(http.MethodGet, "/api/masjidboard/display", nil)
	firstResponse := httptest.NewRecorder()
	s.masjidBoardDisplay(firstResponse, firstRequest)

	if firstResponse.Code != http.StatusOK {
		t.Fatalf("first response status = %d, want 200", firstResponse.Code)
	}
	etag := firstResponse.Header().Get("ETag")
	if etag == "" {
		t.Fatal("first response did not include ETag")
	}
	if got := firstResponse.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", got)
	}

	secondRequest := httptest.NewRequest(http.MethodGet, "/api/masjidboard/display", nil)
	secondRequest.Header.Set("If-None-Match", etag)
	secondResponse := httptest.NewRecorder()
	s.masjidBoardDisplay(secondResponse, secondRequest)

	if secondResponse.Code != http.StatusNotModified {
		t.Fatalf("conditional response status = %d, want 304", secondResponse.Code)
	}
	if secondResponse.Body.Len() != 0 {
		t.Fatalf("conditional response body = %q, want empty", secondResponse.Body.String())
	}
}

func TestMasjidBoardDisplayRejectsWrongMethod(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodPost, "/api/masjidboard/display", nil)
	res := httptest.NewRecorder()
	s.masjidBoardDisplay(res, req)
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}
