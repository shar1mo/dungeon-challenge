package app

import (
	"reflect"
	"testing"

	"github.com/shar1mo/dungeon-challenge/internal/config"
	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

func TestProcessRegisterPlayer(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
	}

	assertLines(t, got, want)

	player := processor.players[1]
	if player == nil || !player.Registered {
		t.Fatalf("expected player to be registered")
	}
}

func TestProcessUnregisteredPlayerEntersAndGetsDisqualified(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 2, ""),
	})

	want := []string{
		"[14:00:00] Player [1] disqualified",
	}

	assertLines(t, got, want)

	player := processor.players[1]
	if player.State != domain.StateDisqual || !player.Finished {
		t.Fatalf("expected player to be disqualified and finished")
	}
}

func TestProcessEventAfterDisqualificationIgnored(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 2, ""),
		event("14:01:00", 1, 1, ""),
	})

	want := []string{
		"[14:00:00] Player [1] disqualified",
	}

	assertLines(t, got, want)
}

func TestProcessDoubleRegisterImpossible(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 1, ""),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:01:00] Player [1] makes impossible move [1]",
	}

	assertLines(t, got, want)
}

func TestProcessEnterAfterRegisterWorks(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:10:00", 1, 2, ""),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:10:00] Player [1] entered the dungeon",
	}

	assertLines(t, got, want)

	player := processor.players[1]
	if !player.Started {
		t.Fatalf("expected player to be started")
	}
}

func TestProcessLeaveBeforeEnterImpossible(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:05:00", 1, 8, ""),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:05:00] Player [1] makes impossible move [8]",
	}

	assertLines(t, got, want)

	player := processor.players[1]
	if player.Finished {
		t.Fatalf("leave before enter should not finish player")
	}
}

func TestProcessCannotContinueSetsDisqual(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:10:00", 1, 9, "too tired to continue"),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:10:00] Player [1] cannot continue due to [too tired to continue]",
	}

	assertLines(t, got, want)

	player := processor.players[1]
	if player.State != domain.StateDisqual || !player.Finished {
		t.Fatalf("expected player to be disqualified and finished")
	}
}

func newTestProcessor(t *testing.T) *Processor {
	t.Helper()

	cfg := config.Config{
		Floors:   2,
		Monsters: 2,
		OpenAt:   "14:05:00",
		Duration: 2,
	}

	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("new processor: %v", err)
	}

	return processor
}

func event(clock string, playerID int, eventID int, extra string) domain.Event {
	seconds, err := domain.ParseClock(clock)
	if err != nil {
		panic(err)
	}

	return domain.Event{
		TimeSeconds: seconds,
		PlayerID:    playerID,
		EventID:     eventID,
		Extra:       extra,
	}
}

func assertLines(t *testing.T, got []string, want []string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("\nwant: %#v\ngot:  %#v", want, got)
	}
}
