package player

import (
	"errors"
	"reflect"
	"testing"
)

func TestMixerDevice(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"alsa/plughw:CARD=UACDemoV10,DEV=0", "hw:CARD=UACDemoV10", true},
		{"alsa/hw:CARD=UACDemoV10,DEV=0", "hw:CARD=UACDemoV10", true},
		{"alsa/plughw:1,0", "hw:1", true},
		{"alsa/default", "", false},
	}

	for _, test := range tests {
		got, ok := mixerDevice(test.input)
		if got != test.want || ok != test.ok {
			t.Fatalf("mixerDevice(%q) = %q, %v; want %q, %v", test.input, got, ok, test.want, test.ok)
		}
	}
}

func TestALSAVolumeGet(t *testing.T) {
	calls := make([][]string, 0)
	mixer := &ALSAVolume{run: func(name string, args ...string) ([]byte, error) {
		calls = append(calls, append([]string{name}, args...))
		return []byte("  Front Left: Playback 75 [75%] [0.00dB]\n"), nil
	}}

	volume, supported, err := mixer.Get("alsa/plughw:CARD=USB,DEV=0")
	if err != nil || !supported || volume != 75 {
		t.Fatalf("Get() = %d, %v, %v; want 75, true, nil", volume, supported, err)
	}
	wantCall := []string{"amixer", "-D", "hw:CARD=USB", "sget", "Master"}
	if !reflect.DeepEqual(calls[0], wantCall) {
		t.Fatalf("amixer call = %#v; want %#v", calls[0], wantCall)
	}
}

func TestALSAVolumeSet(t *testing.T) {
	calls := make([][]string, 0)
	mixer := &ALSAVolume{run: func(name string, args ...string) ([]byte, error) {
		calls = append(calls, append([]string{name}, args...))
		if args[len(args)-1] == "sget" || args[len(args)-2] == "sget" {
			return []byte("  Front Left: Playback 40 [40%] [0.00dB]\n"), nil
		}
		return nil, nil
	}}

	supported, err := mixer.Set("alsa/plughw:CARD=USB,DEV=0", 40)
	if err != nil || !supported {
		t.Fatalf("Set() = %v, %v; want true, nil", supported, err)
	}
	if len(calls) != 2 {
		t.Fatalf("expected two amixer calls, got %d", len(calls))
	}
	wantSet := []string{"amixer", "-D", "hw:CARD=USB", "sset", "Master", "40%"}
	if !reflect.DeepEqual(calls[1], wantSet) {
		t.Fatalf("set call = %#v; want %#v", calls[1], wantSet)
	}
}

func TestALSAVolumeUnsupported(t *testing.T) {
	mixer := &ALSAVolume{run: func(name string, args ...string) ([]byte, error) {
		return []byte("amixer: Unable to find simple control 'Master'"), errors.New("control missing")
	}}

	_, supported, err := mixer.Get("alsa/plughw:CARD=USB,DEV=0")
	if err == nil || supported {
		t.Fatalf("expected unsupported hardware error, got supported=%v err=%v", supported, err)
	}
}
