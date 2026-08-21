package api

import "testing"

func TestCurrentInstalledComponentsDefaultsToBoth(t *testing.T) {
	t.Setenv("MASJIDPI_COMPONENTS", "")
	got := currentInstalledComponents()
	if !got.Listen || !got.Board {
		t.Fatalf("components = %+v, want both enabled", got)
	}
}

func TestCurrentInstalledComponentsProfiles(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		listen bool
		board  bool
	}{
		{name: "listen", value: "listen", listen: true},
		{name: "board", value: "board", board: true},
		{name: "both", value: "listen,board", listen: true, board: true},
		{name: "both alias", value: "both", listen: true, board: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MASJIDPI_COMPONENTS", tt.value)
			got := currentInstalledComponents()
			if got.Listen != tt.listen || got.Board != tt.board {
				t.Fatalf("components = %+v, want listen=%v board=%v", got, tt.listen, tt.board)
			}
		})
	}
}
