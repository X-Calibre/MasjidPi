package masjidboardlive

import (
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

func TestApplyCoreUpcomingSalaahChangeFallbackUsesRenderedHTML(t *testing.T) {
	result := CoreResult{}
	page := []byte(`<section id="nextChangeHeadingC">
		<h1 id="nextChangeHeading">Next <u>S</u>alaah Change</h1>
		<h3 id="fajrNextDate">15 Sep</h3><h3 id="fajrNextTime">05:15</h3>
		<h3 id="asrNextDate">06 Oct</h3><h3 id="asrNextTime">17:00</h3>
		<h3 id="eshaNextDate">06 Oct</h3><h3 id="eshaNextTime">20:00</h3>
	</section>`)

	applyCoreUpcomingSalaahChangeFallback(&result, page)

	if len(result.Board.Notices) != 3 {
		t.Fatalf("notices = %+v, want three Salaah changes", result.Board.Notices)
	}
	want := []struct {
		prayer string
		date   string
		time   string
	}{
		{prayer: "Fajr", date: "15 Sep", time: "05:15"},
		{prayer: "Asr", date: "06 Oct", time: "17:00"},
		{prayer: "Esha", date: "06 Oct", time: "20:00"},
	}
	for index, expected := range want {
		notice := result.Board.Notices[index]
		if notice.Type != model.NoticeTypeSalaahChange ||
			notice.Fields["prayer"] != expected.prayer ||
			notice.Fields["effective_date"] != expected.date ||
			notice.Fields["new_time"] != expected.time {
			t.Fatalf("notice %d = %+v, want %+v", index, notice, expected)
		}
	}
}

func TestApplyCoreUpcomingSalaahChangeFallbackSkipsMissingValues(t *testing.T) {
	result := CoreResult{}
	page := []byte(`<h3 id="fajrNextDate">&ndash;</h3><h3 id="fajrNextTime">&nbsp;</h3>
		<h3 id="asrNextDate">06 Oct</h3><h3 id="asrNextTime">17:00</h3>
		<h3 id="eshaNextDate">00:00</h3><h3 id="eshaNextTime">00:00</h3>`)

	applyCoreUpcomingSalaahChangeFallback(&result, page)

	if len(result.Board.Notices) != 1 || result.Board.Notices[0].Fields["prayer"] != "Asr" {
		t.Fatalf("notices = %+v, want only Asr", result.Board.Notices)
	}
}

func TestApplyCoreUpcomingSalaahChangeFallbackDoesNotDuplicateStructuredNotice(t *testing.T) {
	result := CoreResult{Board: model.Board{Notices: []model.Notice{{
		Type: model.NoticeTypeSalaahChange,
		Fields: map[string]string{
			"prayer":         "Fajr",
			"effective_date": "15 Sep",
			"new_time":       "05:15",
		},
	}}}}
	page := []byte(`<h3 id="asrNextDate">06 Oct</h3><h3 id="asrNextTime">17:00</h3>`)

	applyCoreUpcomingSalaahChangeFallback(&result, page)

	if len(result.Board.Notices) != 1 {
		t.Fatalf("notices = %+v, existing structured values were duplicated", result.Board.Notices)
	}
}
