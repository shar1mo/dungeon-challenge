package domain

import "testing"

func TestParseClockValid(t *testing.T) {
	seconds, err := ParseClock("14:05:09")
	if err != nil {
		t.Fatalf("parse clock: %v", err)
	}

	expected := 14*3600 + 5*60 + 9
	if seconds != expected {
		t.Fatalf("expected %d, got %d", expected, seconds)
	}
}

func TestParseClockInvalid(t *testing.T) {
	values := []string{
		"14:5:09",
		"14:05",
		"24:00:00",
		"10:60:00",
		"10:00:60",
		"aa:00:00",
		"10:bb:00",
		"10:00:cc",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			if _, err := ParseClock(value); err == nil {
				t.Fatalf("expected error for %q", value)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	got := FormatDuration(2*3600 + 5*60 + 9)
	want := "02:05:09"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormatClock(t *testing.T) {
	got := FormatClock(25*3600 + 5*60 + 9)
	want := "01:05:09"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
