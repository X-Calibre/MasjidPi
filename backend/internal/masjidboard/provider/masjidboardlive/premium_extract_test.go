package masjidboardlive

import (
	"strings"
	"testing"
)

func premiumRowsJSON() string {
	rows := make([]string, 29)
	for i := range rows {
		rows[i] = "[]"
	}
	rows[0] = `["quoted ] value", "escaped \" value"]`
	return "[" + strings.Join(rows, ",") + "]"
}

func TestExtractPremiumPageData(t *testing.T) {
	html := []byte(`<script>let boardId = "opaque-id";</script>` +
		`<script>let theInfo = ` + premiumRowsJSON() + `;</script>`)

	boardID, rows, err := ExtractPremiumPageData(html)
	if err != nil {
		t.Fatalf("ExtractPremiumPageData() error = %v", err)
	}
	if boardID != "opaque-id" {
		t.Fatalf("boardID = %q", boardID)
	}
	if len(rows) != 29 {
		t.Fatalf("rows = %d, want 29", len(rows))
	}
}

func TestExtractPremiumPageDataRejectsMissingBoardID(t *testing.T) {
	if _, _, err := ExtractPremiumPageData([]byte(`<script>let theInfo = ` + premiumRowsJSON() + `;</script>`)); err == nil {
		t.Fatal("expected missing board ID error")
	}
}

func TestExtractPremiumPageDataRejectsWrongRowCount(t *testing.T) {
	html := []byte(`<script>let boardId = "opaque-id"; let theInfo = [[],[]];</script>`)
	if _, _, err := ExtractPremiumPageData(html); err == nil {
		t.Fatal("expected wrong row count error")
	}
}
