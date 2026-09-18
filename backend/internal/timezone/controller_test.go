package timezone

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type fakeRunner struct {
	name string
	err  error
}

func (r *fakeRunner) SetTimezone(
	_ context.Context,
	name string,
) error {
	r.name = name
	return r.err
}

func testController(t *testing.T, runner *fakeRunner) *Controller {
	t.Helper()

	table := filepath.Join(t.TempDir(), "zone.tab")
	data := "" +
		"ZA\t-2615+02800\tAfrica/Johannesburg\n" +
		"US\t+404251-0740023\tAmerica/New_York\tEastern\n" +
		"US\t+340308-1181434\tAmerica/Los_Angeles\tPacific\n"
	if err := os.WriteFile(table, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	return newController(
		filepath.Join(t.TempDir(), "timezone.json"),
		table,
		runner,
	)
}

func TestZonesReturnsCountryTimezones(t *testing.T) {
	controller := testController(t, &fakeRunner{})

	zones, err := controller.Zones("South Africa")
	if err != nil {
		t.Fatal(err)
	}
	if len(zones) != 1 || zones[0].Name != "Africa/Johannesburg" {
		t.Fatalf("zones = %+v", zones)
	}

	zones, err = controller.Zones("United States")
	if err != nil {
		t.Fatal(err)
	}
	if len(zones) != 2 ||
		zones[0].Name != "America/Los_Angeles" ||
		zones[1].Name != "America/New_York" {
		t.Fatalf("zones = %+v", zones)
	}
}

func TestZonesRejectsUnknownCountry(t *testing.T) {
	controller := testController(t, &fakeRunner{})

	if _, err := controller.Zones("Unknown"); err == nil {
		t.Fatal("Zones() expected an error")
	}
}

func TestSetAppliesAndPersistsTimezone(t *testing.T) {
	runner := &fakeRunner{}
	controller := testController(t, runner)

	if err := controller.Set(
		context.Background(),
		"Africa/Johannesburg",
	); err != nil {
		t.Fatal(err)
	}
	if runner.name != "Africa/Johannesburg" {
		t.Fatalf("runner name = %q", runner.name)
	}

	data, err := os.ReadFile(controller.statePath)
	if err != nil {
		t.Fatal(err)
	}
	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatal(err)
	}
	if state.Name != "Africa/Johannesburg" {
		t.Fatalf("persisted name = %q", state.Name)
	}
}

func TestSetRejectsUnknownTimezone(t *testing.T) {
	runner := &fakeRunner{}
	controller := testController(t, runner)

	err := controller.Set(context.Background(), "Africa/Not_A_Zone")
	if err == nil {
		t.Fatal("Set() expected an error")
	}
	if runner.name != "" {
		t.Fatalf("runner was called with %q", runner.name)
	}
}

func TestSetDoesNotPersistCommandFailure(t *testing.T) {
	runner := &fakeRunner{err: errors.New("command failed")}
	controller := testController(t, runner)

	err := controller.Set(
		context.Background(),
		"Africa/Johannesburg",
	)
	if err == nil {
		t.Fatal("Set() expected an error")
	}
	if _, statErr := os.Stat(controller.statePath); !errors.Is(
		statErr,
		os.ErrNotExist,
	) {
		t.Fatalf("state file exists after failed command: %v", statErr)
	}
}

func TestRestoreAppliesPersistedTimezone(t *testing.T) {
	runner := &fakeRunner{}
	controller := testController(t, runner)

	if err := os.WriteFile(
		controller.statePath,
		[]byte(`{"name":"Africa/Johannesburg"}`),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	restored, err := controller.Restore(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !restored || runner.name != "Africa/Johannesburg" {
		t.Fatalf("restored=%v runner=%q", restored, runner.name)
	}
}

func TestRestoreMissingStateIsNoOp(t *testing.T) {
	controller := testController(t, &fakeRunner{})

	restored, err := controller.Restore(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if restored {
		t.Fatal("Restore() = true for missing state")
	}
}
