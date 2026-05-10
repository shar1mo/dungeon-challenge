package domain

import (
	"fmt"
	"strconv"
	"strings"
)

const secondsInDay = 24 * 60 * 60

func ParseClock(openat string) (int, error) {
	parts := strings.Split(openat, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid clock format: %q", openat)
	}

	for _, part := range parts {
		if len(part) != 2 {
			return 0, fmt.Errorf("invalid clock format: %q", openat)
		}
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %w", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	if hours < 0 || hours > 23 {
		return 0, fmt.Errorf("hourse out of range: %d", hours)
	}

	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("minutes out of range: %d", minutes)
	}

	if seconds < 0 || seconds > 59 {
		return 0, fmt.Errorf("seconds out of range: %d", seconds)
	}

	return hours*3600 + minutes*60 + seconds, nil
}

func FormatDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}

	hours := seconds / 3600
	seconds %= 3600

	minutes := seconds / 60
	seconds %= 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func FormatClock(seconds int) string {
	seconds %= secondsInDay
	if seconds < 0 {
		seconds += secondsInDay
	}

	return FormatDuration(seconds)
}
