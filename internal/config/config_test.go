package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	data := []byte(`{
		"Floors": 2,
		"Monsters": 3,
		"OpenAt": "14:05:00",
		"Duration": 2
	}`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Floors != 2 {
		t.Fatalf("expected Floors 2, got %d", cfg.Floors)
	}

	if cfg.Monsters != 3 {
		t.Fatalf("expected Monsters 3, got %d", cfg.Monsters)
	}

	if cfg.OpenAt != "14:05:00" {
		t.Fatalf("expected OpenAt 14:05:00, got %q", cfg.OpenAt)
	}

	if cfg.Duration != 2 {
		t.Fatalf("expected Duration 2, got %d", cfg.Duration)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		ok   bool
	}{
		{
			name: "valid",
			cfg: Config{
				Floors:   2,
				Monsters: 2,
				OpenAt:   "14:05:00",
				Duration: 2,
			},
			ok: true,
		},
		{
			name: "invalid floors",
			cfg: Config{
				Floors:   0,
				Monsters: 2,
				OpenAt:   "14:05:00",
				Duration: 2,
			},
			ok: false,
		},
		{
			name: "invalid monsters",
			cfg: Config{
				Floors:   2,
				Monsters: 0,
				OpenAt:   "14:05:00",
				Duration: 2,
			},
			ok: false,
		},
		{
			name: "invalid open at",
			cfg: Config{
				Floors:   2,
				Monsters: 2,
				OpenAt:   "14:5:00",
				Duration: 2,
			},
			ok: false,
		},
		{
			name: "invalid duration",
			cfg: Config{
				Floors:   2,
				Monsters: 2,
				OpenAt:   "14:05:00",
				Duration: 0,
			},
			ok: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()

			if tt.ok && err != nil {
				t.Fatalf("expected valid config, got error: %v", err)
			}

			if !tt.ok && err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}
