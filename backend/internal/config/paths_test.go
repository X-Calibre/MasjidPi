package config

import (
	"path/filepath"
	"testing"
)

func TestNewPathsUpdateStateUsesDevelopmentDataRoot(
	t *testing.T,
) {
	root := t.TempDir()
	paths := NewPaths(root)

	want := filepath.Join(
		root,
		"backend",
		"data",
		"update_state.json",
	)

	if paths.UpdateState != want {
		t.Fatalf(
			"UpdateState = %q, want %q",
			paths.UpdateState,
			want,
		)
	}
	if filepath.Dir(paths.UpdateState) != paths.DataRoot {
		t.Fatalf(
			"UpdateState directory = %q, want DataRoot %q",
			filepath.Dir(paths.UpdateState),
			paths.DataRoot,
		)
	}
}

func TestRuntimePathsUpdateStateUsesPersistentDataRoot(
	t *testing.T,
) {
	t.Setenv("MASJIDPI_HOME", "/opt/masjidpi")

	paths, err := RuntimePaths()
	if err != nil {
		t.Fatal(err)
	}

	const want = "/var/lib/masjidpi/update_state.json"

	if paths.UpdateState != want {
		t.Fatalf(
			"UpdateState = %q, want %q",
			paths.UpdateState,
			want,
		)
	}
	if filepath.Dir(paths.UpdateState) != paths.DataRoot {
		t.Fatalf(
			"UpdateState directory = %q, want DataRoot %q",
			filepath.Dir(paths.UpdateState),
			paths.DataRoot,
		)
	}
}
