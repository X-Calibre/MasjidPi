package playback

import "testing"

type fakeVolumePersistence struct {
	saves []struct {
		device string
		volume int
	}
}

func (f *fakeVolumePersistence) Load(string) (int, bool, error) { return 0, false, nil }

func (f *fakeVolumePersistence) Save(device string, volume int) error {
	f.saves = append(f.saves, struct {
		device string
		volume int
	}{device: device, volume: volume})
	return nil
}

func TestSetVolumeTransientDoesNotPersist(t *testing.T) {
	player := &fakePlayer{}
	store := &fakeVolumePersistence{}
	manager := New(player, Config{})
	manager.SetVolumePersistence(store)

	if err := manager.SetVolumeTransient(42); err != nil {
		t.Fatalf("set transient volume: %v", err)
	}

	if len(store.saves) != 0 {
		t.Fatalf("persistent saves = %d, want 0", len(store.saves))
	}
	if got := manager.Status().Volume; got != 42 {
		t.Fatalf("manager volume = %d, want 42", got)
	}
	if got := player.volumeCallsSnapshot(); len(got) != 1 || got[0] != 42 {
		t.Fatalf("player volume calls = %v, want [42]", got)
	}
}

func TestPersistVolumeWritesCurrentVolumeOnce(t *testing.T) {
	player := &fakePlayer{}
	store := &fakeVolumePersistence{}
	manager := New(player, Config{})
	manager.SetVolumePersistence(store)

	if err := manager.SetVolumeTransient(42); err != nil {
		t.Fatalf("set transient volume: %v", err)
	}
	if err := manager.PersistVolume(); err != nil {
		t.Fatalf("persist volume: %v", err)
	}

	if len(store.saves) != 1 {
		t.Fatalf("persistent saves = %d, want 1", len(store.saves))
	}
	if store.saves[0].device != "auto" || store.saves[0].volume != 42 {
		t.Fatalf("saved volume = %+v, want device auto and volume 42", store.saves[0])
	}
}
