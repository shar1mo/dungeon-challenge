package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shar1mo/dungeon-challenge/internal/config"
	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

func main() {
	configPath := flag.String("config", "testdata/config.json", "path to config file")
	eventsPath := flag.String("events", "testdata/events", "path to events file")
	flag.Parse()

	if _, err := config.Load(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if _, err := domain.ReadEvents(*eventsPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read events: %v\n", err)
		os.Exit(1)
	}
}
