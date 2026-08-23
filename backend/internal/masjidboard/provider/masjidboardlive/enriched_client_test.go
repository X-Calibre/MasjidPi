package masjidboardlive

import (
	"context"
	"errors"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
)

type enrichmentProvider struct {
	board model.Board
	err   error
	calls int
}

func (p *enrichmentProvider) Fetch(context.Context) (model.Board, error) {
	p.calls++
	return p.board, p.err
}

var _ provider.Provider = EnrichedClient{}

func TestEnrichedClientMergesOptionalPremiumContent(t *testing.T) {
	coreProvider := &enrichmentProvider{board: model.Board{
		Identity: model.BoardIdentity{ID: "core-id", Name: "Core Name", TimeZone: "GMT+02:00"},
	}}
	premiumProvider := &enrichmentProvider{board: model.Board{
		Identity:      model.BoardIdentity{ID: "premium-id", Name: "Premium Name", TimeZone: "GMT+02"},
		Announcements: []model.Announcement{{Title: "Announcement", Content: "Community update"}},
		Notices:       []model.Notice{{Type: model.NoticeTypeFuneral, Title: "Funeral Notice"}},
	}}

	board, err := (EnrichedClient{Core: coreProvider, Premium: premiumProvider}).Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if board.Identity.ID != "core-id" || board.Identity.Name != "Core Name" {
		t.Fatalf("Core identity was replaced: %+v", board.Identity)
	}
	if len(board.Announcements) != 1 || board.Announcements[0].Title != "Announcement" {
		t.Fatalf("Announcements = %+v", board.Announcements)
	}
	if len(board.Notices) != 1 || board.Notices[0].Type != model.NoticeTypeFuneral {
		t.Fatalf("Notices = %+v", board.Notices)
	}
}

func TestEnrichedClientReturnsCoreWhenPremiumFails(t *testing.T) {
	coreProvider := &enrichmentProvider{board: model.Board{
		Identity:      model.BoardIdentity{ID: "core-id", Name: "Core Name"},
		Announcements: []model.Announcement{{Title: "Previously normalised Core content"}},
	}}
	premiumProvider := &enrichmentProvider{err: errors.New("Premium unavailable")}

	board, err := (EnrichedClient{Core: coreProvider, Premium: premiumProvider}).Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if board.Identity.ID != "core-id" || len(board.Announcements) != 1 {
		t.Fatalf("Core board was not preserved: %+v", board)
	}
}

func TestEnrichedClientReturnsCoreFailure(t *testing.T) {
	coreProvider := &enrichmentProvider{err: errors.New("Core unavailable")}
	premiumProvider := &enrichmentProvider{board: model.Board{Identity: model.BoardIdentity{ID: "premium-id"}}}

	if _, err := (EnrichedClient{Core: coreProvider, Premium: premiumProvider}).Fetch(context.Background()); err == nil {
		t.Fatal("Fetch() expected Core error")
	}
	if premiumProvider.calls != 0 {
		t.Fatalf("Premium calls = %d, want 0 after Core failure", premiumProvider.calls)
	}
}

func TestEnrichedClientRejectsMissingCoreProvider(t *testing.T) {
	if _, err := (EnrichedClient{}).Fetch(context.Background()); err == nil {
		t.Fatal("Fetch() expected missing Core provider error")
	}
}
