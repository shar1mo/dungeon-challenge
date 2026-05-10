package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if cfg.Floors <= 0 {
		return fmt.Errorf("floors must be positive")
	}

	if cfg.Monsters <= 0 {
		return fmt.Errorf("monsters must be positive")
	}

	if cfg.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	if _, err := cfg.OpenAtSeconds(); err != nil {
		return fmt.Errorf("open at: %w", err)
	}

	return nil
}

func (cfg Config) OpenAtSeconds() (int, error) {
	return domain.ParseClock(cfg.OpenAt)
}

func (cfg Config) CloseAtSeconds() (int, error) {
	openAt, err := cfg.OpenAtSeconds()
	if err != nil {
		return 0, err
	}

	return openAt + cfg.Duration*3600, nil
}