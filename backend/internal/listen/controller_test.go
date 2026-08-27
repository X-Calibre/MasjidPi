package listen

import (
	"fmt"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type fakeAvailability struct {
	available map[string]bool
}

func (f *fakeAvailability) Status(mount string) (bool, bool, time.Time) {
	available, ok := f.available[mount]
	return available, ok, time.Now()
}

func (f *fakeAvailability) Events() <-chan string { return nil }

type activation struct {
	id     string
	volume int
}

type fakeOutput struct {
	activations []activation
	volumes     []int
	stops       int
}

func (f *fakeOutput) Activate(selected stream.Stream, volume int) error {
	f.activations = append(f.activations, activation{id: selected.ID, volume: volume})
	return nil
}

func (f *fakeOutput) SetSoftwareVolume(volume int) error {
	f.volumes = append(f.volumes, volume)
	return nil
}

func (f *fakeOutput) Stop() error {
	f.stops++
	return nil
}

func TestControllerUsesRadioUntilMasjidComesOnline(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	masjid := &stream.Stream{ID: "masjid-a", Name: "Masjid A", URL: "https://example.test/masjid"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, Name: "Radio A", URL: "https://example.test/radio"}

	if err := controller.SelectMasjid(masjid); err != nil {
		t.Fatal(err)
	}
	if err := controller.SelectRadio(radio); err != nil {
		t.Fatal(err)
	}
	controller.Listen()
	controller.step()

	if got := controller.Status().ActiveSource; got != ActiveRadio {
		t.Fatalf("active source = %q, want %q", got, ActiveRadio)
	}
	if len(output.activations) != 1 || output.activations[0].id != radio.ID || output.activations[0].volume != 70 {
		t.Fatalf("radio activation = %#v", output.activations)
	}

	availability.available[masjid.ID] = true
	controller.step()
	if got := controller.Status().ActiveSource; got != ActiveMasjid {
		t.Fatalf("active source = %q, want %q", got, ActiveMasjid)
	}
	if len(output.activations) != 2 || output.activations[1].id != masjid.ID || output.activations[1].volume != 100 {
		t.Fatalf("masjid activation = %#v", output.activations)
	}
}

func TestControllerDelaysRadioAfterActiveMasjidGoesOffline(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": true}}
	output := &fakeOutput{}
	controller := New(availability, output)
	controller.radioResumeDelay = 20 * time.Millisecond
	masjid := &stream.Stream{ID: "masjid-a", URL: "https://example.test/masjid"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()

	availability.available[masjid.ID] = false
	controller.step()
	status := controller.Status()
	if status.ActiveSource != ActiveNone || !status.RadioResumePending {
		t.Fatalf("status after masjid offline = %+v, want no active source with pending radio resume", status)
	}

	time.Sleep(25 * time.Millisecond)
	controller.step()
	status = controller.Status()
	if status.ActiveSource != ActiveRadio || status.RadioResumePending {
		t.Fatalf("status after delay = %+v, want radio active and no pending delay", status)
	}
}

func TestControllerCanOverridePendingRadioDelay(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": true}}
	output := &fakeOutput{}
	controller := New(availability, output)
	controller.radioResumeDelay = time.Minute
	masjid := &stream.Stream{ID: "masjid-a", URL: "https://example.test/masjid"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()

	availability.available[masjid.ID] = false
	controller.step()
	if err := controller.ResumeRadioNow(); err != nil {
		t.Fatalf("ResumeRadioNow: %v", err)
	}
	controller.step()
	status := controller.Status()
	if status.ActiveSource != ActiveRadio || status.RadioResumePending {
		t.Fatalf("status after manual override = %+v, want radio active immediately", status)
	}
}

func TestControllerCancelsRadioDelayWhenMasjidReturns(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": true}}
	output := &fakeOutput{}
	controller := New(availability, output)
	controller.radioResumeDelay = time.Minute
	masjid := &stream.Stream{ID: "masjid-a", URL: "https://example.test/masjid"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()

	availability.available[masjid.ID] = false
	controller.step()
	availability.available[masjid.ID] = true
	controller.step()

	status := controller.Status()
	if status.ActiveSource != ActiveMasjid || status.RadioResumePending {
		t.Fatalf("status after masjid return = %+v, want immediate masjid and cancelled delay", status)
	}
}

func TestControllerRadioScheduleSilencesOutsideWindow(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectRadio(radio)

	now := time.Now()
	start := now.Add(2 * time.Minute)
	stop := now.Add(3 * time.Minute)
	startText := fmt.Sprintf("%02d:%02d", start.Hour(), start.Minute())
	stopText := fmt.Sprintf("%02d:%02d", stop.Hour(), stop.Minute())
	if err := controller.SetRadioSchedule(true, startText, stopText); err != nil {
		t.Fatal(err)
	}
	controller.Listen()
	controller.step()

	status := controller.Status()
	if status.ActiveSource != ActiveNone || status.RadioScheduleAllowsNow {
		t.Fatalf("status outside schedule = %+v, want silence", status)
	}
}

func TestRadioWindowAllowsDaytimeAndOvernightRanges(t *testing.T) {
	day := time.Date(2026, 8, 27, 12, 0, 0, 0, time.Local)
	if !radioWindowAllows(day, "06:00", "22:00") {
		t.Fatal("12:00 should be inside 06:00-22:00")
	}
	if radioWindowAllows(time.Date(2026, 8, 27, 23, 0, 0, 0, time.Local), "06:00", "22:00") {
		t.Fatal("23:00 should be outside 06:00-22:00")
	}
	if !radioWindowAllows(time.Date(2026, 8, 27, 23, 0, 0, 0, time.Local), "22:00", "02:00") {
		t.Fatal("23:00 should be inside overnight 22:00-02:00")
	}
	if !radioWindowAllows(time.Date(2026, 8, 28, 1, 0, 0, 0, time.Local), "22:00", "02:00") {
		t.Fatal("01:00 should be inside overnight 22:00-02:00")
	}
}

func TestControllerAppliesIndependentActiveSourceVolumes(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	masjid := &stream.Stream{ID: "masjid-a", URL: "https://example.test/masjid"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()

	if err := controller.SetRadioVolume(55); err != nil {
		t.Fatal(err)
	}
	if len(output.volumes) != 1 || output.volumes[0] != 55 {
		t.Fatalf("software volume updates = %#v, want [55]", output.volumes)
	}

	if err := controller.SetMasjidVolume(90); err != nil {
		t.Fatal(err)
	}
	if len(output.volumes) != 1 {
		t.Fatalf("inactive masjid volume changed live output: %#v", output.volumes)
	}

	availability.available[masjid.ID] = true
	controller.step()
	last := output.activations[len(output.activations)-1]
	if last.id != masjid.ID || last.volume != 90 {
		t.Fatalf("masjid activation = %#v, want id=%q volume=90", last, masjid.ID)
	}
}

func TestControllerStopsAllAudioWhenListeningDisabled(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()
	controller.Stop()
	controller.step()

	status := controller.Status()
	if status.ActiveSource != ActiveNone {
		t.Fatalf("active source = %q, want none", status.ActiveSource)
	}
	if output.stops != 1 {
		t.Fatalf("stop count = %d, want 1", output.stops)
	}
}

func TestControllerRejectsWrongSourceKinds(t *testing.T) {
	controller := New(nil, &fakeOutput{})
	radio := &stream.Stream{ID: "radio", Kind: stream.KindRadio}
	masjid := &stream.Stream{ID: "masjid"}

	if err := controller.SelectMasjid(radio); err == nil {
		t.Fatal("expected radio to be rejected as primary masjid")
	}
	if err := controller.SelectRadio(masjid); err == nil {
		t.Fatal("expected masjid to be rejected as secondary radio")
	}
}
