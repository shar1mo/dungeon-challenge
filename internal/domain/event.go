package domain

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Event struct {
	TimeSeconds int
	PlayerID    int
	EventID     int
	Extra       string
}

func ParseEventLine(line string) (Event, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return Event{}, fmt.Errorf("invalid event line: %q", line)
	}

	timeToken := fields[0]
	if len(timeToken) != len("[HH:MM:SS]") || timeToken[0] != '[' || timeToken[len(timeToken)-1] != ']' {
		return Event{}, fmt.Errorf("invalid event time format: %q", timeToken)
	}

	timeSeconds, err := ParseClock(timeToken[1 : len(timeToken)-1])
	if err != nil {
		return Event{}, err
	}

	playerID, err := strconv.Atoi(fields[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid player id: %w", err)
	}

	eventID, err := strconv.Atoi(fields[2])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event id: %w", err)
	}

	event := Event{
		TimeSeconds: timeSeconds,
		PlayerID:    playerID,
		EventID:     eventID,
	}

	if len(fields) > 3 {
		event.Extra = strings.Join(fields[3:], " ")
	}

	return event, nil
}

func ReadEvents(path string) ([]Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Event

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := ParseEventLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
