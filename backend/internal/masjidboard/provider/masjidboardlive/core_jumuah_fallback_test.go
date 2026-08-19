package masjidboardlive

import (
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

func TestApplyCoreJumuahHeadingFallbackUsesRenderedHTML(t *testing.T) {
	result := CoreResult{Board: model.Board{PrayerTimes: model.PrayerTimes{Jumuah: []model.JumuahService{{
		Events: []model.JumuahEvent{
			{Time: &model.ClockTime{Hour: 12, Minute: 25}},
			{Time: &model.ClockTime{Hour: 12, Minute: 55}},
		},
	}}}}}

	page := []byte(`<table id="jumuahTable"><tr>
		<td><h2 id="jumuahHead1">Adhan</h2></td>
		<td><h2 id="jumuahHead2">Sunan</h2></td>
		<td><h2 id="jumuahHead3">Khu<u>t</u>bah</h2></td>
	</tr><tr>
		<td><h3 id="jumuahTime1">12:25</h3></td>
		<td><h3 id="jumuahTime2">&nbsp;</h3></td>
		<td><h3 id="jumuahTime3">12:55</h3></td>
	</tr></table>`)

	applyCoreJumuahHeadingFallback(&result, page)

	service := result.Board.PrayerTimes.Jumuah[0]
	if len(service.Events) != 2 {
		t.Fatalf("events = %d, want 2", len(service.Events))
	}
	if service.Events[0].Code != "0" || service.Events[0].Heading != "Adhan" {
		t.Fatalf("event 1 = code %q heading %q, want 0/Adhan", service.Events[0].Code, service.Events[0].Heading)
	}
	if service.Events[1].Code != "6" || service.Events[1].Heading != "Khutbah" {
		t.Fatalf("event 2 = code %q heading %q, want 6/Khutbah", service.Events[1].Code, service.Events[1].Heading)
	}
	if service.Adhan == nil || service.Adhan.Hour != 12 || service.Adhan.Minute != 25 {
		t.Fatalf("Adhan = %#v, want 12:25", service.Adhan)
	}
}

func TestApplyCoreJumuahHeadingFallbackDoesNotOverwriteStructuredHeading(t *testing.T) {
	result := CoreResult{Board: model.Board{PrayerTimes: model.PrayerTimes{Jumuah: []model.JumuahService{{
		Events: []model.JumuahEvent{{Code: "1", Heading: "Lecture", Time: &model.ClockTime{Hour: 12, Minute: 40}}},
	}}}}}

	page := []byte(`<h2 id="jumuahHead1">Adhan</h2><h3 id="jumuahTime1">12:40</h3>`)
	applyCoreJumuahHeadingFallback(&result, page)

	event := result.Board.PrayerTimes.Jumuah[0].Events[0]
	if event.Code != "1" || event.Heading != "Lecture" {
		t.Fatalf("event = code %q heading %q, structured value was overwritten", event.Code, event.Heading)
	}
}
