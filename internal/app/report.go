package app

import (
	"fmt"
	"sort"

	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

func (p *Processor) Report() []string {
	playerIDs := make([]int, 0, len(p.players))
	for playerID := range p.players {
		playerIDs = append(playerIDs, playerID)
	}

	sort.Ints(playerIDs)

	lines := []string{"Final report:"}

	for _, playerID := range playerIDs {
		player := p.players[playerID]
		lines = append(lines, p.formatReportLine(player))
	}

	return lines
}

func (p *Processor) formatReportLine(player *domain.Player) string {
	totalTime := 0
	if player.Started {
		totalTime = player.FinishedAt - player.EnteredAt
	}

	avgFloorTime := 0
	if player.ClearedFloors > 0 {
		avgFloorTime = player.TotalFloorClearTime / player.ClearedFloors
	}

	bossTime := 0
	if player.BossKilled {
		bossTime = player.BossKillDuration
	}

	state := player.State
	if state == "" {
		state = domain.StateFail
	}

	return fmt.Sprintf(
		"[%s] %d [%s, %s, %s] HP:%d",
		state,
		player.ID,
		domain.FormatDuration(totalTime),
		domain.FormatDuration(avgFloorTime),
		domain.FormatDuration(bossTime),
		player.Health,
	)
}
