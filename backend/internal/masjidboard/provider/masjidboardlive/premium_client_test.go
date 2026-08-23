package masjidboardlive

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
)

var _ provider.Provider = PremiumClient{}

func TestPremiumClientFetchAt(t *testing.T) {
	rows := loadCapturedRows(t)
	payload, err := json.Marshal(rows)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	var gotMid string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMid = r.URL.Query().Get("mid")
		_, _ = w.Write([]byte(`<script>let boardId = "opaque-id"; let theInfo = ` + string(payload) + `;</script>`))
	}))
	defer server.Close()

	client := PremiumClient{HTTPClient: server.Client(), Endpoint: server.URL, Mid: "azaadville-darul-uloom"}
	board, err := client.FetchAt(context.Background(), time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchAt() error = %v", err)
	}
	if gotMid != "azaadville-darul-uloom" {
		t.Fatalf("mid = %q", gotMid)
	}
	if board.Identity.ID != "azaadville-darul-uloom" {
		t.Fatalf("Identity.ID = %q", board.Identity.ID)
	}
}

func TestPremiumClientRejectsPageWithoutPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>not a board</html>`))
	}))
	defer server.Close()

	client := PremiumClient{HTTPClient: server.Client(), Endpoint: server.URL, Mid: "missing"}
	if _, err := client.FetchAt(context.Background(), time.Now()); err == nil {
		t.Fatal("expected missing payload error")
	}
}

func TestPremiumClientRejectsMissingMid(t *testing.T) {
	if _, err := (PremiumClient{}).FetchAt(context.Background(), time.Now()); err == nil {
		t.Fatal("expected missing mid error")
	}
}
