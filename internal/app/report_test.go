package app

import (
	"testing"

	"github.com/shar1mo/dungeon-challenge/internal/config"
	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

func TestReportSuccess(t *testing.T) {
	processor := newTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
		event("14:02:00", 1, 3, ""),
		event("14:03:00", 1, 3, ""),
		event("14:04:00", 1, 4, ""),
		event("14:05:00", 1, 3, ""),
		event("14:06:00", 1, 3, ""),
		event("14:07:00", 1, 6, ""),
		event("14:10:00", 1, 7, ""),
		event("14:11:00", 1, 8, ""),
	})

	want := []string{
		"Final report:",
		"[SUCCESS] 1 [00:10:00, 00:02:00, 00:03:00] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestReportFailOnDeath(t *testing.T) {
	processor := newTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
		event("14:02:00", 1, 11, "100"),
	})

	want := []string{
		"Final report:",
		"[FAIL] 1 [00:01:00, 00:00:00, 00:00:00] HP:0",
	}

	assertLines(t, processor.Report(), want)
}

func TestReportDisqual(t *testing.T) {
	processor := newTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 9, "bad luck"),
	})

	want := []string{
		"Final report:",
		"[DISQUAL] 1 [00:00:00, 00:00:00, 00:00:00] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestReportAvgFloorTime(t *testing.T) {
	processor := newTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
		event("14:03:00", 1, 3, ""),
		event("14:05:00", 1, 3, ""),
		event("14:06:00", 1, 4, ""),
		event("14:07:00", 1, 3, ""),
		event("14:09:00", 1, 3, ""),
		event("14:10:00", 1, 8, ""),
	})

	want := []string{
		"Final report:",
		"[FAIL] 1 [00:09:00, 00:03:30, 00:00:00] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestReportBossTime(t *testing.T) {
	processor := newTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
		event("14:02:00", 1, 3, ""),
		event("14:03:00", 1, 3, ""),
		event("14:04:00", 1, 4, ""),
		event("14:05:00", 1, 3, ""),
		event("14:06:00", 1, 3, ""),
		event("14:10:00", 1, 6, ""),
		event("14:17:00", 1, 7, ""),
		event("14:18:00", 1, 8, ""),
	})

	want := []string{
		"Final report:",
		"[SUCCESS] 1 [00:17:00, 00:02:00, 00:07:00] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestReportSorting(t *testing.T) {
	processor := newTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 3, 1, ""),
		event("14:00:00", 1, 1, ""),
		event("14:00:00", 2, 1, ""),
	})

	want := []string{
		"Final report:",
		"[FAIL] 1 [00:00:00, 00:00:00, 00:00:00] HP:100",
		"[FAIL] 2 [00:00:00, 00:00:00, 00:00:00] HP:100",
		"[FAIL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestDungeonCloseFail(t *testing.T) {
	processor := newClosingTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
	})

	want := []string{
		"Final report:",
		"[FAIL] 1 [00:04:00, 00:00:00, 00:00:00] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestDungeonCloseSuccess(t *testing.T) {
	processor := newClosingTestProcessor(t)

	processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
		event("14:02:00", 1, 3, ""),
		event("14:03:00", 1, 3, ""),
		event("14:04:00", 1, 4, ""),
		event("14:04:10", 1, 3, ""),
		event("14:04:20", 1, 3, ""),
		event("14:04:30", 1, 6, ""),
		event("14:04:50", 1, 7, ""),
	})

	want := []string{
		"Final report:",
		"[SUCCESS] 1 [00:04:00, 00:01:10, 00:00:20] HP:100",
	}

	assertLines(t, processor.Report(), want)
}

func TestGoldenSampleInputOutput(t *testing.T) {
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

	events := []domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:00:00", 2, 1, ""),
		event("14:10:00", 2, 2, ""),
		event("14:10:00", 3, 2, ""),
		event("14:11:00", 2, 5, ""),
		event("14:12:00", 3, 3, ""),
		event("14:14:00", 2, 3, ""),
		event("14:27:00", 2, 11, "60"),
		event("14:29:00", 2, 11, "50"),
		event("14:40:00", 1, 2, ""),
		event("14:41:00", 1, 3, ""),
		event("14:44:00", 1, 11, "50"),
		event("14:45:00", 1, 3, ""),
		event("14:48:00", 1, 4, ""),
		event("14:48:00", 1, 6, ""),
		event("14:49:00", 1, 11, "25"),
		event("14:49:02", 1, 10, "80"),
		event("14:50:00", 1, 11, "65"),
		event("14:59:00", 1, 7, ""),
		event("15:04:00", 1, 8, ""),
	}

	processor.Process(events)

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:00:00] Player [2] registered",
		"[14:10:00] Player [2] entered the dungeon",
		"[14:10:00] Player [3] is disqualified",
		"[14:11:00] Player [2] makes imposible move [5]",
		"[14:14:00] Player [2] killed the monster",
		"[14:27:00] Player [2] recieved [60] of damage",
		"[14:29:00] Player [2] recieved [50] of damage",
		"[14:29:00] Player [2] is dead",
		"[14:40:00] Player [1] entered the dungeon",
		"[14:41:00] Player [1] killed the monster",
		"[14:44:00] Player [1] recieved [50] of damage",
		"[14:45:00] Player [1] killed the monster",
		"[14:48:00] Player [1] went to the next floor",
		"[14:48:00] Player [1] entered the boss's floor",
		"[14:49:00] Player [1] recieved [25] of damage",
		"[14:49:02] Player [1] has restored [80] of health",
		"[14:50:00] Player [1] recieved [65] of damage",
		"[14:59:00] Player [1] killed the boss",
		"[15:04:00] Player [1] left the dungeon",
		"Final report:",
		"[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35",
		"[FAIL] 2 [00:19:00, 00:00:00, 00:00:00] HP:0",
		"[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100",
	}

	assertLines(t, processor.OutputWithReport(), want)
}

func newClosingTestProcessor(t *testing.T) *Processor {
	t.Helper()

	cfg := config.Config{
		Floors:   3,
		Monsters: 2,
		OpenAt:   "14:00:00",
		Duration: 0,
	}

	cfg.Duration = 1

	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("new processor: %v", err)
	}

	processor.closeAt = mustClock("14:05:00")

	return processor
}

func mustClock(value string) int {
	seconds, err := domain.ParseClock(value)
	if err != nil {
		panic(err)
	}

	return seconds
}

func TestProcessEnterBeforeOpenAtImpossible(t *testing.T) {
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

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:01:00] Player [1] makes imposible move [2]",
	}

	assertLines(t, got, want)
}

func TestProcessNegativeDamageImpossible(t *testing.T) {
	processor := newTestProcessor(t)

	got := processor.Process([]domain.Event{
		event("14:00:00", 1, 1, ""),
		event("14:01:00", 1, 2, ""),
		event("14:02:00", 1, 11, "-10"),
	})

	want := []string{
		"[14:00:00] Player [1] registered",
		"[14:01:00] Player [1] entered the dungeon",
		"[14:02:00] Player [1] makes imposible move [11]",
	}

	assertLines(t, got, want)
}
