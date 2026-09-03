package dailycontent

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestClientLive is opt-in because it depends on the external MasjidBoard
// Live service. Set MASJIDBOARD_DAILY_CONTENT_LIVE_TEST=1 when explicitly
// validating the current upstream JavaScript contract.
func TestClientLive(t *testing.T) {
	if os.Getenv("MASJIDBOARD_DAILY_CONTENT_LIVE_TEST") != "1" {
		t.Skip("set MASJIDBOARD_DAILY_CONTENT_LIVE_TEST=1 to run the live daily-content test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	content, err := (Client{}).Fetch(ctx)
	if err != nil {
		t.Fatalf("Client.Fetch() error = %v", err)
	}
	if !content.Valid() {
		t.Fatalf("incomplete live daily content: %+v", content)
	}
	t.Logf("daily content parsed: date=%s language=%s", content.ContentDate, content.Language)
}
