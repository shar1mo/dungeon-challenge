package domain

import "testing"

func TestParseEventLineWithoutExtra(t *testing.T) {
	event, err := ParseEventLine("[14:00:00] 1 1")
	if err != nil {
		t.Fatalf("parse event line: %v", err)
	}

	if event.TimeSeconds != 14*3600 {
		t.Fatalf("unexpected time: %d", event.TimeSeconds)
	}

	if event.PlayerID != 1 {
		t.Fatalf("unexpected player id: %d", event.PlayerID)
	}

	if event.EventID != 1 {
		t.Fatalf("unexpected event id: %d", event.EventID)
	}

	if event.Extra != "" {
		t.Fatalf("expected empty extra, got %q", event.Extra)
	}
}

func TestParseEventLineWithNumericExtra(t *testing.T) {
	event, err := ParseEventLine("[14:27:00] 2 11 60")
	if err != nil {
		t.Fatalf("parse event line: %v", err)
	}

	if event.Extra != "60" {
		t.Fatalf("expected extra %q, got %q", "60", event.Extra)
	}
}

func TestParseEventLineWithMultiWordExtra(t *testing.T) {
	event, err := ParseEventLine("[14:29:00] 2 9 too tired to continue")
	if err != nil {
		t.Fatalf("parse event line: %v", err)
	}

	want := "too tired to continue"
	if event.Extra != want {
		t.Fatalf("expected extra %q, got %q", want, event.Extra)
	}
}

func TestParseEventLineInvalid(t *testing.T) {
	values := []string{
		"",
		"[14:00:00] 1",
		"14:00:00 1 1",
		"[14:00] 1 1",
		"[14:00:00] abc 1",
		"[14:00:00] 1 abc",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			if _, err := ParseEventLine(value); err == nil {
				t.Fatalf("expected error for %q", value)
			}
		})
	}
}
