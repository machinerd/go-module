package log

import "testing"

func TestLevelFromSlogRoundTrip(t *testing.T) {
	for _, l := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError} {
		if got := levelFromSlog(l.slogLevel()); got != l {
			t.Errorf("levelFromSlog(%v.slogLevel()) = %v, want %v", l, got, l)
		}
	}
}

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input  string
		want   Level
		wantOk bool
	}{
		{"debug", LevelDebug, true},
		{"INFO", LevelInfo, true},
		{"Warn", LevelWarn, true},
		{"warning", LevelWarn, true},
		{"error", LevelError, true},
		{"trace", 0, false},
		{"", 0, false},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := ParseLevel(tc.input)
			if ok != tc.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOk)
			}
			if ok && got != tc.want {
				t.Errorf("level = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	cases := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
	}

	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}
