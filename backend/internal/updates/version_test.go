package updates

import "testing"

func TestParseStableVersion(t *testing.T) {
	tests := []struct {
		value string
		want  StableVersion
	}{
		{
			value: "v1.6.0",
			want: StableVersion{
				Major: 1,
				Minor: 6,
				Patch: 0,
			},
		},
		{
			value: "v10.20.30",
			want: StableVersion{
				Major: 10,
				Minor: 20,
				Patch: 30,
			},
		},
		{
			value: "v0.0.0",
			want:  StableVersion{},
		},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			got, err := ParseStableVersion(test.value)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf(
					"ParseStableVersion(%q) = %+v, want %+v",
					test.value,
					got,
					test.want,
				)
			}
		})
	}
}

func TestParseStableVersionRejectsNonStableTags(t *testing.T) {
	values := []string{
		"",
		"1.6.0",
		"v1.6",
		"v1.6.0-rc.1",
		"v1.6.0-image",
		"v1.6.0-lab.1",
		"v01.6.0",
		"v1.06.0",
		"v1.6.00",
		"v1.6.0+build",
		"latest",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			if _, err := ParseStableVersion(value); err == nil {
				t.Fatalf(
					"ParseStableVersion(%q) error = nil",
					value,
				)
			}
		})
	}
}

func TestStableVersionCompare(t *testing.T) {
	tests := []struct {
		name  string
		left  StableVersion
		right StableVersion
		want  int
	}{
		{
			name:  "equal",
			left:  StableVersion{1, 6, 0},
			right: StableVersion{1, 6, 0},
			want:  0,
		},
		{
			name:  "newer major",
			left:  StableVersion{2, 0, 0},
			right: StableVersion{1, 99, 99},
			want:  1,
		},
		{
			name:  "older minor",
			left:  StableVersion{1, 5, 99},
			right: StableVersion{1, 6, 0},
			want:  -1,
		},
		{
			name:  "newer patch",
			left:  StableVersion{1, 6, 10},
			right: StableVersion{1, 6, 9},
			want:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.left.Compare(test.right); got != test.want {
				t.Fatalf(
					"Compare() = %d, want %d",
					got,
					test.want,
				)
			}
		})
	}
}
