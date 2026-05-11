package app

import (
	"fmt"

	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

func formatEventLine(timeSeconds int, message string) string {
	return fmt.Sprintf("[%s] %s", domain.FormatClock(timeSeconds), message)
}

func formatRegistered(timeSeconds int, playerID int) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] registered", playerID))
}

func formatEnteredDungeon(timeSeconds int, playerID int) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] entered the dungeon", playerID))
}

func formatLeftDungeon(timeSeconds int, playerID int) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] left the dungeon", playerID))
}

func formatCannotContinue(timeSeconds int, playerID int, reason string) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] cannot continue due to [%s]", playerID, reason))
}

func formatDisqualified(timeSeconds int, playerID int) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] disqualified", playerID))
}

func formatDead(timeSeconds int, playerID int) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] is dead", playerID))
}

func formatImpossibleMove(timeSeconds int, playerID int, eventID int) string {
	return formatEventLine(timeSeconds, fmt.Sprintf("Player [%d] makes impossible move [%d]", playerID, eventID))
}
