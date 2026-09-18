package updates

import "testing"

func TestIsNewerThanInstalled(t *testing.T) {
	tests := []struct {
		name      string
		candidate StableVersion
		current   string
		want      bool
	}{
		{
			name:      "newer stable",
			candidate: StableVersion{1, 7, 0},
			current:   "v1.6.0",
			want:      true,
		},
		{
			name:      "same stable",
			candidate: StableVersion{1, 6, 0},
			current:   "v1.6.0-image",
			want:      false,
		},
		{
			name:      "older stable",
			candidate: StableVersion{1, 5, 9},
			current:   "v1.6.0",
			want:      false,
		},
		{
			name:      "stable supersedes same RC",
			candidate: StableVersion{1, 6, 0},
			current:   "v1.6.0-rc.4-image",
			want:      true,
		},
		{
			name:      "older than installed RC core",
			candidate: StableVersion{1, 6, 0},
			current:   "v1.7.0-rc.1-image",
			want:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := IsNewerThanInstalled(
				test.candidate,
				test.current,
			)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf(
					"IsNewerThanInstalled() = %v, want %v",
					got,
					test.want,
				)
			}
		})
	}
}

func TestParseInstalledVersionRejectsUnsafeIdentity(
	t *testing.T,
) {
	values := []string{
		"",
		"v1.6.0-lab.1",
		"v1.6.0-35-gabcdef0",
		"v1.6.0-rc.0",
		"1.6.0",
		"v01.6.0",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			if _, err := parseInstalledVersion(value); err == nil {
				t.Fatalf(
					"parseInstalledVersion(%q) error = nil",
					value,
				)
			}
		})
	}
}
