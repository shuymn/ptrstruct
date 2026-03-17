package ptrstruct

import "testing"

func TestMatchNolint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		text        string
		honorAll    bool
		wantMatched bool
	}{
		{name: "exact ptrstruct", text: "//nolint:ptrstruct", wantMatched: true},
		{name: "ptrstruct with space", text: "// nolint:ptrstruct", wantMatched: true},
		{name: "ptrstruct with reason", text: "//nolint:ptrstruct // legacy API", wantMatched: true},
		{name: "nolint all honored", text: "//nolint:all", honorAll: true, wantMatched: true},
		{name: "nolint all not honored", text: "//nolint:all", honorAll: false, wantMatched: false},
		{name: "nolint all with space", text: "// nolint:all", honorAll: true, wantMatched: true},
		{name: "other linter", text: "//nolint:govet", wantMatched: false},
		{name: "multiple linters including ptrstruct", text: "//nolint:govet,ptrstruct", wantMatched: true},
		{name: "regular comment", text: "// this is a comment", wantMatched: false},
		{name: "empty comment", text: "//", wantMatched: false},
		{name: "nolint bare", text: "//nolint", wantMatched: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := matchNolint(tt.text, tt.honorAll)
			if got != tt.wantMatched {
				t.Errorf("matchNolint(%q, %v) = %v, want %v", tt.text, tt.honorAll, got, tt.wantMatched)
			}
		})
	}
}
