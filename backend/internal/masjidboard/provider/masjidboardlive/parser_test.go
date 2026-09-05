package masjidboardlive

import (
	"embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

//go:embed testdata/azaadville-darul-uloom-core.json
var capturedFixture embed.FS

func loadCapturedRows(t *testing.T) []json.RawMessage {
	t.Helper()
	data, err := capturedFixture.ReadFile("testdata/azaadville-darul-uloom-core.json")
	if err != nil {
		t.Fatalf("read captured fixture: %v", err)
	}
	var rows []json.RawMessage
	if err := json.Unmarshal(data, &rows); err != nil {
		t.Fatalf("decode captured fixture: %v", err)
	}
	return rows
}

func TestParseCapturedMasjidBoardLiveCore(t *testing.T) {
	rows := loadCapturedRows(t)
	if len(rows) < 29 {
		t.Fatalf("fixture rows = %d, want at least 29", len(rows))
	}

	now := time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC)
	board, err := Parse(rows, "1Zpg5LKfd_ZoEQsA0rsyWNBrUgY6QVaHnGdPfuKHF24A", now)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if board.Identity.ID != "azaadville-darul-uloom" {
		t.Fatalf("Identity.ID = %q", board.Identity.ID)
	}
	if board.Identity.Name != "Azaadville" {
		t.Fatalf("Identity.Name = %q", board.Identity.Name)
	}
	if board.Identity.AlternateName != "Madrasah Arabia Islamia" {
		t.Fatalf("Identity.AlternateName = %q", board.Identity.AlternateName)
	}
	if board.Identity.TimeZone != "GMT+02" {
		t.Fatalf("Identity.TimeZone = %q", board.Identity.TimeZone)
	}
	if got := board.DateContext.GregorianDate.Format("2006-01-02"); got != "2026-09-11" {
		t.Fatalf("GregorianDate = %q", got)
	}

	assertClock := func(name string, got *model.ClockTime, hour, minute int) {
		t.Helper()
		if got == nil {
			t.Fatalf("%s is nil", name)
		}
		if got.Hour != hour || got.Minute != minute {
			t.Fatalf("%s = %02d:%02d, want %02d:%02d", name, got.Hour, got.Minute, hour, minute)
		}
	}

	assertClock("Fajr Adhan", board.PrayerTimes.Fajr.Adhan, 5, 30)
	assertClock("Fajr Jamaah", board.PrayerTimes.Fajr.Jamaah, 5, 50)
	assertClock("Dhuhr Adhan", board.PrayerTimes.Dhuhr.Adhan, 12, 30)
	assertClock("Dhuhr Jamaah", board.PrayerTimes.Dhuhr.Jamaah, 12, 45)
	assertClock("Asr Adhan", board.PrayerTimes.Asr.Adhan, 16, 15)
	assertClock("Asr Jamaah", board.PrayerTimes.Asr.Jamaah, 16, 30)
	assertClock("Maghrib Adhan", board.PrayerTimes.Maghrib.Adhan, 17, 52)
	if board.PrayerTimes.Maghrib.Jamaah != nil {
		t.Fatalf("Maghrib Jamaah = %+v, want nil", board.PrayerTimes.Maghrib.Jamaah)
	}
	assertClock("Esha Adhan", board.PrayerTimes.Esha.Adhan, 19, 15)
	assertClock("Esha Jamaah", board.PrayerTimes.Esha.Jamaah, 19, 30)
	if board.PrayerTimes.SpecialDhuhr == nil {
		t.Fatal("Special Dhuhr is nil")
	}
	assertClock("Special Dhuhr", board.PrayerTimes.SpecialDhuhr.Time, 12, 45)
	if board.PrayerTimes.SpecialDhuhr.Label != "(Everyday)" {
		t.Fatalf("Special Dhuhr label = %q", board.PrayerTimes.SpecialDhuhr.Label)
	}

	if len(board.PrayerTimes.Jumuah) != 1 {
		t.Fatalf("Jumuah services = %d, want 1", len(board.PrayerTimes.Jumuah))
	}
	jumuah := board.PrayerTimes.Jumuah[0]
	assertClock("Jumuah Adhan", jumuah.Adhan, 12, 45)
	assertClock("Jumuah Jamaah", jumuah.Jamaah, 12, 55)
	assertClock("Jumuah Islamic Adhan", jumuah.IslamicAdhan, 18, 56)
	assertClock("Jumuah Islamic Jamaah", jumuah.IslamicJamaah, 19, 6)
	if jumuah.AlternateAdhan != nil || jumuah.AlternateJamaah != nil {
		t.Fatal("newly parsed Jumuah data must not use deprecated alternate fields")
	}
	if jumuah.Khateeb != "Sunnats after Adhān" {
		t.Fatalf("Jumuah Khateeb = %q", jumuah.Khateeb)
	}
	if len(jumuah.Events) != 3 {
		t.Fatalf("Jumuah events = %d, want 3", len(jumuah.Events))
	}
	if jumuah.Events[0].Code != "1" || jumuah.Events[0].Heading != "Lecture" {
		t.Fatalf("Jumuah event 1 = %+v", jumuah.Events[0])
	}
	assertClock("Jumuah event 1", jumuah.Events[0].Time, 12, 15)
	if jumuah.Events[1].Code != "0" || jumuah.Events[1].Heading != "Adhān" {
		t.Fatalf("Jumuah event 2 = %+v", jumuah.Events[1])
	}
	assertClock("Jumuah event 2", jumuah.Events[1].Time, 12, 45)
	if jumuah.Events[2].Code != "6" || jumuah.Events[2].Heading != "Khutbah" {
		t.Fatalf("Jumuah event 3 = %+v", jumuah.Events[2])
	}
	assertClock("Jumuah event 3", jumuah.Events[2].Time, 12, 55)
	assertClock("Jumuah effective Salaah", jumuah.EffectiveSalaah(), 12, 55)

	if board.Astronomical == nil {
		t.Fatal("Astronomical is nil")
	}
	assertClock("Suhur", board.Astronomical.Suhur, 5, 19)
	assertClock("FajrStart", board.Astronomical.FajrStart, 5, 19)
	assertClock("Sunrise", board.Astronomical.Sunrise, 6, 37)
	assertClock("Ishraaq", board.Astronomical.Ishraaq, 6, 52)
	assertClock("Duha", board.Astronomical.Duha, 9, 25)
	assertClock("AsrShafii", board.Astronomical.AsrShafii, 15, 25)
	assertClock("AsrHanafi", board.Astronomical.AsrHanafi, 16, 13)
	assertClock("Sunset", board.Astronomical.Sunset, 17, 49)
	assertClock("EshaStart", board.Astronomical.EshaStart, 19, 7)
}

func TestParseJumuahDoesNotInferSalaahFromKhutbah(t *testing.T) {
	row := json.RawMessage(`["Adhān","12:25","Sunan","12:55","Khutbah","13:00","Ml M Bhamjee","12:25","","18:35","19:10","0,3,6"]`)
	service, err := parseJumuahRow(row)
	if err != nil {
		t.Fatalf("parseJumuahRow() error = %v", err)
	}
	if service.Jamaah != nil {
		t.Fatal("Jumuah Jamaah should be absent")
	}
	if got := service.EffectiveSalaah(); got != nil {
		t.Fatalf("EffectiveSalaah() = %+v, want nil without explicit Jamaah", got)
	}
	if len(service.Events) != 3 || service.Events[2].Heading != "Khutbah" || service.Events[2].Time == nil || service.Events[2].Time.Hour != 13 {
		t.Fatalf("Jumuah events = %+v", service.Events)
	}
}

func TestParseJumuahAllowsEmptyRow(t *testing.T) {
	row := json.RawMessage(`["","","","","","","","","","","","#N/A"]`)
	service, err := parseJumuahRow(row)
	if err != nil {
		t.Fatalf("parseJumuahRow() error = %v", err)
	}
	if service != nil {
		t.Fatalf("service = %+v, want nil", service)
	}
}

func TestCorePlaceholderIsAbsent(t *testing.T) {
	if !isAbsent("~~~~") {
		t.Fatal("expected ~~~~ to be treated as absent")
	}
}

func TestParseRejectsIncomplete29RowResponse(t *testing.T) {
	rows := loadCapturedRows(t)[:28]
	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error for fewer than 29 rows")
	}
}

func TestParseRejectsMissingTimezone(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowClock] = json.RawMessage(`[]`)
	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error for missing timezone")
	}
}

func TestParseRejectsMissingCorePrayer(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowSalah] = json.RawMessage(`["","", "12:30","12:45", "16:15","16:30", "17:52","", "19:15","19:30"]`)
	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error when a core prayer has no usable time")
	}
}

func TestParseAllowsMissingOptionalPrayerValue(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowSalah] = json.RawMessage(`["05:30","", "12:30","12:45", "16:15","16:30", "17:52","", "19:15","19:30"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if board.PrayerTimes.Fajr.Adhan == nil || board.PrayerTimes.Fajr.Jamaah != nil {
		t.Fatal("expected Fajr Adhan only")
	}
	if board.PrayerTimes.Maghrib.Adhan == nil || board.PrayerTimes.Maghrib.Jamaah != nil {
		t.Fatal("expected Maghrib Adhan only")
	}
}

func TestParsePromotesAlternateMasjidNameWhenPrimaryIsEmpty(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowMasjid] = json.RawMessage(`["","Darul Uloom Zakariyya","https://masjidboardlive.com","2","7200000"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if board.Identity.Name != "Darul Uloom Zakariyya" {
		t.Fatalf("Identity.Name = %q", board.Identity.Name)
	}
	if board.Identity.AlternateName != "" {
		t.Fatalf("Identity.AlternateName = %q, want empty", board.Identity.AlternateName)
	}
}

func TestParseCommunityContentFromPremiumRows(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowAnnouncement] = json.RawMessage(`["Urgent Announcement","<b>Tonight after Maghrib</b>","Display","Hidden","Do not show","Hide"]`)
	rows[rowAnnouncement2] = json.RawMessage(`[]`)
	rows[rowNikah] = json.RawMessage(`["Ahmad","Son of","Yusuf","Daughter of","Ismail","Sat, 29 Aug 2026","after Asr","1787961600000","Display","FALSE","Maryam"]`)
	rows[rowFuneral] = json.RawMessage(`["Abdullah","father of Zaid","10 Main Road","14:00","Central Cemetery","Masjid Hall","14:30","Display"]`)
	rows[rowEid] = json.RawMessage(`["Mon, 25 May 2026","Sports Ground","1 Field Road","07:00","07:30","Display"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(board.Announcements) != 1 {
		t.Fatalf("Announcements = %+v", board.Announcements)
	}
	if board.Announcements[0].Title != "Urgent Announcement" || board.Announcements[0].Content != "<b>Tonight after Maghrib</b>" {
		t.Fatalf("Announcement = %+v", board.Announcements[0])
	}
	if len(board.Notices) != 3 {
		t.Fatalf("Notices = %+v", board.Notices)
	}
	if board.Notices[0].Type != model.NoticeTypeNikah || board.Notices[0].Fields["bride"] != "Maryam" {
		t.Fatalf("Nikah notice = %+v", board.Notices[0])
	}
	if board.Notices[1].Type != model.NoticeTypeFuneral || board.Notices[1].Fields["salaah_time"] != "14:30" {
		t.Fatalf("Funeral notice = %+v", board.Notices[1])
	}
	if board.Notices[2].Type != model.NoticeTypeEid || board.Notices[2].Fields["venue"] != "Sports Ground" {
		t.Fatalf("Eid notice = %+v", board.Notices[2])
	}
}

func TestParseCommunityContentIgnoresHiddenRows(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowAnnouncement] = json.RawMessage(`["Announcement","Hidden content","Hide"]`)
	rows[rowAnnouncement2] = json.RawMessage(`[]`)
	rows[rowNikah] = json.RawMessage(`["Name","Relation","Name","Relation","Name","Date","Time","0","Hide"]`)
	rows[rowFuneral] = json.RawMessage(`["Name","Relation","Address","Pickup","Cemetery","Venue","Time","Hide"]`)
	rows[rowEid] = json.RawMessage(`["Date","Venue","Address","Lecture","Salaah","Hide"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(board.Announcements) != 0 || len(board.Notices) != 0 {
		t.Fatalf("hidden community content was retained: announcements=%+v notices=%+v", board.Announcements, board.Notices)
	}
}

func TestParseAdditionalCommunityCardContent(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowUpcoming] = json.RawMessage(`["–","–","0","–","–","0","–","–","0","TRUE","24 Aug","05:50","1787529600000","1 Sep","17:00","1788220800000","4 Sep","19:45","1788480000000"]`)
	rows[rowClock] = json.RawMessage(`["23 Aug, 05:27","23 Aug, 18:34","12 hrs","273.10","7.25","24 Aug, 18:37","24 Aug, 19:32","36 hrs","272.18","19.48","23 August 2026","24 August 2026","test-board","","Johannesburg","GMT+02","en","Islamic Time","Display Islamic Time","","","TRUE"]`)
	rows[rowTaleem] = json.RawMessage(`["Wednesday 11:15–12:15","26 August 2026","Resident's home","10 Main Road","Display","","","","","Hide"]`)
	rows[rowWellWishes] = json.RawMessage(`["Please make du'a for the unwell","","","","","","","","","","Display"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 8, 23, 9, 0, 0, 0, time.FixedZone("SAST", 2*60*60)))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(board.Notices) != 4 {
		t.Fatalf("Notices = %+v", board.Notices)
	}
	if board.Notices[0].Type != model.NoticeTypeSalaahChange || board.Notices[0].Fields["new_time"] != "05:50" {
		t.Fatalf("Salaah change = %+v", board.Notices[0])
	}
	if board.Notices[3].Type != model.NoticeTypeWellWish {
		t.Fatalf("Well wishes = %+v", board.Notices[3])
	}
	if len(board.Programmes) != 1 || board.Programmes[0].Title != "Taleem Programme" {
		t.Fatalf("Programmes = %+v", board.Programmes)
	}
	if board.NewMoon == nil || board.NewMoon.Fields["visibility_date"] != "24 August 2026" {
		t.Fatalf("NewMoon = %+v", board.NewMoon)
	}
}

func TestParseAdditionalCommunityCardContentHonoursVisibility(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowUpcoming] = json.RawMessage(`["–","–","0","–","–","0","–","–","0","FALSE","–","–","0","–","–","0","–","–","0"]`)
	rows[rowTaleem] = json.RawMessage(`["Programme","Date","","","Hide","","","","","Hide"]`)
	rows[rowWellWishes] = json.RawMessage(`["Hidden message","","","","","","","","","","Hide"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(board.Programmes) != 0 {
		t.Fatalf("hidden programmes = %+v", board.Programmes)
	}
	for _, notice := range board.Notices {
		if notice.Type == model.NoticeTypeSalaahChange || notice.Type == model.NoticeTypeWellWish {
			t.Fatalf("hidden notice retained: %+v", notice)
		}
	}
}

func TestParseAnnouncementCategoriesAreConservative(t *testing.T) {
	raw := json.RawMessage(`[
		"Salah Times Change","New times from Monday","Display",
		"Class Time Changes","Morning class at 08:00","Display",
		"Weekly Programs","Tafseer after Esha","Display",
		"Taraweeh 2026","One juz nightly","Display",
		"Urgent access notice","Use the side entrance","Display"
	]`)
	announcements := parseAnnouncementRows(raw)
	want := []string{"salaah_change_announcement", "class_time_change", "weekly_programme", "ramadan_programme", "urgent_announcement"}
	if len(announcements) != len(want) {
		t.Fatalf("announcements = %+v", announcements)
	}
	for index, category := range want {
		if announcements[index].Category != category {
			t.Fatalf("announcement %d category=%q want %q", index, announcements[index].Category, category)
		}
	}
}

func TestParseDawahAndContributionCards(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowDawah] = json.RawMessage(`["Daily after Esha","Thursday","Monday","After Asr","After Maghrib","Display","Three-Day Jamaat","Hartbeespoort","4–6 September","Pretoria West","11–13 September","Display"]`)
	rows[rowBanking] = json.RawMessage(`["300000","Masjid Contributions<br>Lillah Only","Example Bank","Example Masjid Trust","123456","000123456","Display",""]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	var dawah, threeDay *model.Notice
	for index := range board.Notices {
		switch board.Notices[index].Type {
		case model.NoticeTypeDawah:
			dawah = &board.Notices[index]
		case model.NoticeTypeThreeDay:
			threeDay = &board.Notices[index]
		}
	}
	if dawah == nil || dawah.Fields["gasht_out_day"] != "Thursday" || dawah.Fields["gasht_in_time"] != "After Maghrib" {
		t.Fatalf("Dawah notice = %+v", dawah)
	}
	if threeDay == nil || threeDay.Fields["first_location"] != "Hartbeespoort" || threeDay.Fields["second_date"] != "11–13 September" {
		t.Fatalf("three-day notice = %+v", threeDay)
	}
	if board.Banking == nil || board.Banking.Content != "Masjid Contributions<br>Lillah Only" || board.Banking.Fields["account_number"] != "000123456" {
		t.Fatalf("Banking = %+v", board.Banking)
	}
}

func TestParseDawahAndContributionCardsHonoursVisibility(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowDawah] = json.RawMessage(`["Daily after Esha","Thursday","Monday","After Asr","After Maghrib","Hide","Three-Day Jamaat","Hartbeespoort","4–6 September","Pretoria West","11–13 September","Hide"]`)
	rows[rowBanking] = json.RawMessage(`["300000","Masjid Contributions","Example Bank","Example Masjid Trust","123456","000123456","Hide",""]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	for _, notice := range board.Notices {
		if notice.Type == model.NoticeTypeDawah || notice.Type == model.NoticeTypeThreeDay {
			t.Fatalf("hidden Dawah notice retained: %+v", notice)
		}
	}
	if board.Banking != nil {
		t.Fatalf("hidden Banking = %+v", board.Banking)
	}
}
