package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shar1mo/dungeon-challenge/internal/app"
	"github.com/shar1mo/dungeon-challenge/internal/config"
	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

func main() {
	configPath := flag.String("config", "testdata/config.json", "path to config file")
	eventsPath := flag.String("events", "testdata/events", "path to events file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	events, err := domain.ReadEvents(*eventsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read events: %v\n", err)
		os.Exit(1)
	}

	processor, err := app.NewProcessor(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create processor: %v\n", err)
		os.Exit(1)
	}

	processor.Process(events)

	for _, line := range processor.OutputWithReport() {
		fmt.Println(line)
	}
}
