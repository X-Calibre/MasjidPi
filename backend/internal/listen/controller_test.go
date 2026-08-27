package listen

import (
	"fmt"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type fakeAvailability struct{ available map[string]bool }

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
func (f *fakeOutput) Stop() error { f.stops++; return nil }

func TestControllerUsesRadioUntilMasjidComesOnline(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	masjid := &stream.Stream{ID: "masjid-a", URL: "https://example.test/masjid"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio, URL: "https://example.test/radio"}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()
	if controller.Status().ActiveSource != ActiveRadio {
		t.Fatal("radio should be active")
	}
	availability.available[masjid.ID] = true
	controller.step()
	if controller.Status().ActiveSource != ActiveMasjid {
		t.Fatal("masjid should interrupt immediately")
	}
}

func TestControllerDelaysRadioAfterActiveMasjidGoesOffline(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": true}}
	output := &fakeOutput{}
	controller := New(availability, output)
	controller.radioResumeDelay = 20 * time.Millisecond
	masjid := &stream.Stream{ID: "masjid-a"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()
	availability.available[masjid.ID] = false
	controller.step()
	if status := controller.Status(); status.ActiveSource != ActiveNone || !status.RadioResumePending {
		t.Fatalf("unexpected hold status: %+v", status)
	}
	time.Sleep(25 * time.Millisecond)
	controller.step()
	if status := controller.Status(); status.ActiveSource != ActiveRadio || status.RadioResumePending {
		t.Fatalf("unexpected post-delay status: %+v", status)
	}
}

func TestPlayNowOverridesQuietTimeAndDelay(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectRadio(radio)
	now := time.Now()
	start := now.Add(2 * time.Minute)
	stop := now.Add(3 * time.Minute)
	_ = controller.SetRadioSchedule(true, fmt.Sprintf("%02d:%02d", start.Hour(), start.Minute()), fmt.Sprintf("%02d:%02d", stop.Hour(), stop.Minute()))
	controller.Listen()
	controller.step()
	if controller.Status().ActiveSource != ActiveNone {
		t.Fatal("schedule should initially silence radio")
	}
	if err := controller.SetRadioMode(RadioModePlayNow); err != nil {
		t.Fatal(err)
	}
	controller.step()
	status := controller.Status()
	if status.ActiveSource != ActiveRadio || status.RadioMode != RadioModePlayNow {
		t.Fatalf("play now did not override quiet time: %+v", status)
	}
}

func TestPlayNowEndsWhenMasjidComesOnline(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	controller := New(availability, &fakeOutput{})
	masjid := &stream.Stream{ID: "masjid-a"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	if err := controller.SetRadioMode(RadioModePlayNow); err != nil {
		t.Fatal(err)
	}
	controller.step()
	availability.available[masjid.ID] = true
	controller.step()
	status := controller.Status()
	if status.ActiveSource != ActiveMasjid || status.RadioMode != RadioModeSchedule {
		t.Fatalf("masjid should cancel play-now override: %+v", status)
	}
}

func TestStoppedRadioRemainsStoppedAcrossMasjidEvents(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	controller := New(availability, &fakeOutput{})
	masjid := &stream.Stream{ID: "masjid-a"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	if err := controller.SetRadioMode(RadioModeStopped); err != nil {
		t.Fatal(err)
	}
	controller.step()
	if controller.Status().ActiveSource != ActiveNone {
		t.Fatal("stopped mode should silence radio")
	}
	availability.available[masjid.ID] = true
	controller.step()
	if controller.Status().ActiveSource != ActiveMasjid {
		t.Fatal("masjid must still play while radio is stopped")
	}
	availability.available[masjid.ID] = false
	controller.step()
	status := controller.Status()
	if status.ActiveSource != ActiveNone || status.RadioMode != RadioModeStopped {
		t.Fatalf("radio stop should persist after masjid event: %+v", status)
	}
	if err := controller.SetRadioMode(RadioModeSchedule); err != nil {
		t.Fatal(err)
	}
	controller.step()
	if controller.Status().ActiveSource != ActiveRadio {
		t.Fatal("play on schedule should re-enable radio")
	}
}

func TestControllerCancelsRadioDelayWhenMasjidReturns(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": true}}
	controller := New(availability, &fakeOutput{})
	controller.radioResumeDelay = time.Minute
	masjid := &stream.Stream{ID: "masjid-a"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
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
		t.Fatalf("unexpected status: %+v", status)
	}
}

func TestControllerRadioScheduleSilencesOutsideWindow(t *testing.T) {
	controller := New(&fakeAvailability{available: map[string]bool{}}, &fakeOutput{})
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectRadio(radio)
	now := time.Now()
	start := now.Add(2 * time.Minute)
	stop := now.Add(3 * time.Minute)
	if err := controller.SetRadioSchedule(true, fmt.Sprintf("%02d:%02d", start.Hour(), start.Minute()), fmt.Sprintf("%02d:%02d", stop.Hour(), stop.Minute())); err != nil {
		t.Fatal(err)
	}
	controller.Listen()
	controller.step()
	status := controller.Status()
	if status.ActiveSource != ActiveNone || status.RadioScheduleAllowsNow {
		t.Fatalf("unexpected schedule status: %+v", status)
	}
}

func TestRadioWindowAllowsDaytimeAndOvernightRanges(t *testing.T) {
	day := time.Date(2026, 8, 27, 12, 0, 0, 0, time.Local)
	if !radioWindowAllows(day, "06:00", "22:00") {
		t.Fatal("12:00 should be inside daytime range")
	}
	if radioWindowAllows(time.Date(2026, 8, 27, 23, 0, 0, 0, time.Local), "06:00", "22:00") {
		t.Fatal("23:00 should be outside daytime range")
	}
	if !radioWindowAllows(time.Date(2026, 8, 27, 23, 0, 0, 0, time.Local), "22:00", "02:00") {
		t.Fatal("23:00 should be inside overnight range")
	}
}

func TestControllerAppliesIndependentActiveSourceVolumes(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	masjid := &stream.Stream{ID: "masjid-a"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()
	if err := controller.SetRadioVolume(55); err != nil {
		t.Fatal(err)
	}
	if len(output.volumes) != 1 || output.volumes[0] != 55 {
		t.Fatalf("unexpected volumes: %#v", output.volumes)
	}
	if err := controller.SetMasjidVolume(90); err != nil {
		t.Fatal(err)
	}
	availability.available[masjid.ID] = true
	controller.step()
	last := output.activations[len(output.activations)-1]
	if last.id != masjid.ID || last.volume != 90 {
		t.Fatalf("unexpected masjid activation: %#v", last)
	}
}

func TestControllerStopsAllAudioWhenListeningDisabled(t *testing.T) {
	output := &fakeOutput{}
	controller := New(&fakeAvailability{available: map[string]bool{}}, output)
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()
	controller.Stop()
	controller.step()
	if controller.Status().ActiveSource != ActiveNone {
		t.Fatal("active source should be none")
	}
}

func TestMasjidPowerOffStopsControllerAndForcesRadioOff(t *testing.T) {
	availability := &fakeAvailability{available: map[string]bool{"masjid-a": false}}
	output := &fakeOutput{}
	controller := New(availability, output)
	masjid := &stream.Stream{ID: "masjid-a"}
	radio := &stream.Stream{ID: "radio-a", Kind: stream.KindRadio}
	_ = controller.SelectMasjid(masjid)
	_ = controller.SelectRadio(radio)
	controller.Listen()
	controller.step()
	if controller.Status().ActiveSource != ActiveRadio {
		t.Fatal("radio should be active before power off")
	}
	controller.SetMasjidEnabled(false)
	controller.step()
	status := controller.Status()
	if status.Listening {
		t.Fatalf("controller should be stopped: %+v", status)
	}
	if status.MasjidEnabled || status.RadioEnabled {
		t.Fatalf("both modules should be off: %+v", status)
	}
	if status.ActiveSource != ActiveNone {
		t.Fatalf("active source should be none: %+v", status)
	}
}

func TestListenCannotStartWhileMasjidPowerIsOff(t *testing.T) {
	controller := New(&fakeAvailability{available: map[string]bool{}}, &fakeOutput{})
	controller.SetMasjidEnabled(false)
	controller.Listen()
	if status := controller.Status(); status.Listening {
		t.Fatalf("controller must remain stopped while masjid is off: %+v", status)
	}
}

func TestControllerRejectsWrongSourceKinds(t *testing.T) {
	controller := New(nil, &fakeOutput{})
	radio := &stream.Stream{ID: "radio", Kind: stream.KindRadio}
	masjid := &stream.Stream{ID: "masjid"}
	if err := controller.SelectMasjid(radio); err == nil {
		t.Fatal("expected radio rejection")
	}
	if err := controller.SelectRadio(masjid); err == nil {
		t.Fatal("expected masjid rejection")
	}
}
