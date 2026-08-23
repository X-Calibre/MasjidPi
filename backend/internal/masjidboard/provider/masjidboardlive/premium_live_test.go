package masjidboardlive

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestPremiumClientLive is opt-in because it depends on the external
// MasjidBoard Live service. Set MASJIDBOARD_LIVE_TEST_MID to a known public
// Premium mid when explicitly validating the current upstream contract.
func TestPremiumClientLive(t *testing.T) {
	mid := os.Getenv("MASJIDBOARD_LIVE_TEST_MID")
	if mid == "" {
		t.Skip("set MASJIDBOARD_LIVE_TEST_MID to run the live Premium test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	board, err := (PremiumClient{Mid: mid}).Fetch(ctx)
	if err != nil {
		t.Fatalf("PremiumClient.Fetch() error = %v", err)
	}
	if board.Identity.ID == "" || board.Identity.Name == "" || board.Identity.TimeZone == "" {
		t.Fatalf("incomplete Premium identity: %+v", board.Identity)
	}
	if board.PrayerTimes.Fajr.Adhan == nil && board.PrayerTimes.Fajr.Jamaah == nil {
		t.Fatal("Premium payload has no usable Fajr time")
	}
	if board.PrayerTimes.Esha.Adhan == nil && board.PrayerTimes.Esha.Jamaah == nil {
		t.Fatal("Premium payload has no usable Esha time")
	}
	t.Logf(
		"Premium board %q parsed: announcements=%d notices=%d programmes=%d media=%d",
		board.Identity.Name,
		len(board.Announcements),
		len(board.Notices),
		len(board.Programmes),
		len(board.Media),
	)
}
