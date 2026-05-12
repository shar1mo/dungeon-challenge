package app

import (
	"strconv"

	"github.com/shar1mo/dungeon-challenge/internal/config"
	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

const (
	eventRegister       = 1
	eventEnterDungeon   = 2
	eventKillMonster    = 3
	eventNextFloor      = 4
	eventPreviousFloor  = 5
	eventEnterBossFloor = 6
	eventKillBoss       = 7
	eventLeaveDungeon   = 8
	eventCannotContinue = 9
	eventRestoreHealth  = 10
	eventReceiveDamage  = 11
)

type Processor struct {
	cfg           config.Config
	players       map[int]*domain.Player
	clearedFloors map[int]map[int]bool
}

func NewProcessor(cfg config.Config) (*Processor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Processor{
		cfg:           cfg,
		players:       make(map[int]*domain.Player),
		clearedFloors: make(map[int]map[int]bool),
	}, nil
}

func (p *Processor) Process(events []domain.Event) []string {
	var output []string

	for _, event := range events {
		player := p.getPlayer(event.PlayerID)

		if player.Finished {
			continue
		}

		if !player.Registered && event.EventID != eventRegister {
			player.State = domain.StateDisqual
			player.Finished = true
			player.FinishedAt = event.TimeSeconds

			output = append(output, formatDisqualified(event.TimeSeconds, player.ID))
			continue
		}

		switch event.EventID {
		case eventRegister:
			output = append(output, p.handleRegister(player, event))

		case eventEnterDungeon:
			output = append(output, p.handleEnterDungeon(player, event))

		case eventKillMonster:
			output = append(output, p.handleKillMonster(player, event))

		case eventNextFloor:
			output = append(output, p.handleNextFloor(player, event))

		case eventPreviousFloor:
			output = append(output, p.handlePreviousFloor(player, event))

		case eventEnterBossFloor:
			output = append(output, p.handleEnterBossFloor(player, event))

		case eventKillBoss:
			output = append(output, p.handleKillBoss(player, event))

		case eventLeaveDungeon:
			output = append(output, p.handleLeaveDungeon(player, event))

		case eventCannotContinue:
			output = append(output, p.handleCannotContinue(player, event))

		case eventRestoreHealth:
			output = append(output, p.handleRestoreHealth(player, event))

		case eventReceiveDamage:
			output = append(output, p.handleReceiveDamage(player, event)...)

		default:
			output = append(output, formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID))
		}
	}

	return output
}

func (p *Processor) getPlayer(id int) *domain.Player {
	player, ok := p.players[id]
	if ok {
		return player
	}

	player = domain.NewPlayer(id)
	p.players[id] = player
	p.clearedFloors[id] = make(map[int]bool)

	return player
}

func (p *Processor) normalFloors() int {
	return p.cfg.Floors - 1
}

func (p *Processor) isActive(player *domain.Player) bool {
	return player.Registered && player.Started && !player.Finished
}

func (p *Processor) isFloorCleared(player *domain.Player, floor int) bool {
	return p.clearedFloors[player.ID][floor]
}

func (p *Processor) markFloorCleared(player *domain.Player, floor int) {
	p.clearedFloors[player.ID][floor] = true
}

func (p *Processor) handleRegister(player *domain.Player, event domain.Event) string {
	if player.Registered {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.Registered = true

	return formatRegistered(event.TimeSeconds, player.ID)
}

func (p *Processor) handleEnterDungeon(player *domain.Player, event domain.Event) string {
	if player.Started {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.Started = true
	player.EnteredAt = event.TimeSeconds
	player.FloorEnterTime = event.TimeSeconds
	player.CurrentFloor = 1

	return formatEnteredDungeon(event.TimeSeconds, player.ID)
}

func (p *Processor) handleKillMonster(player *domain.Player, event domain.Event) string {
	if !p.isActive(player) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.OnBossFloor {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if p.isFloorCleared(player, player.CurrentFloor) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.CurrentFloorKills >= p.cfg.Monsters {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.CurrentFloorKills++

	if player.CurrentFloorKills == p.cfg.Monsters {
		player.ClearedFloors++
		player.TotalFloorClearTime += event.TimeSeconds - player.FloorEnterTime
		p.markFloorCleared(player, player.CurrentFloor)
	}

	return formatKilledMonster(event.TimeSeconds, player.ID)
}

func (p *Processor) handleNextFloor(player *domain.Player, event domain.Event) string {
	if !p.isActive(player) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if !p.isFloorCleared(player, player.CurrentFloor) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.CurrentFloor >= p.normalFloors() {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.CurrentFloor++
	player.FloorEnterTime = event.TimeSeconds

	if p.isFloorCleared(player, player.CurrentFloor) {
		player.CurrentFloorKills = p.cfg.Monsters
	} else {
		player.CurrentFloorKills = 0
	}

	return formatWentNextFloor(event.TimeSeconds, player.ID)
}

func (p *Processor) handlePreviousFloor(player *domain.Player, event domain.Event) string {
	if !p.isActive(player) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.OnBossFloor {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.CurrentFloor <= 1 {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.CurrentFloor--
	player.FloorEnterTime = event.TimeSeconds

	if p.isFloorCleared(player, player.CurrentFloor) {
		player.CurrentFloorKills = p.cfg.Monsters
	} else {
		player.CurrentFloorKills = 0
	}

	return formatWentPreviousFloor(event.TimeSeconds, player.ID)
}

func (p *Processor) handleEnterBossFloor(player *domain.Player, event domain.Event) string {
	if !p.isActive(player) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.ClearedFloors != p.normalFloors() {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.OnBossFloor {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.OnBossFloor = true
	player.BossEnterTime = event.TimeSeconds

	return formatEnteredBossFloor(event.TimeSeconds, player.ID)
}

func (p *Processor) handleKillBoss(player *domain.Player, event domain.Event) string {
	if !p.isActive(player) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if !player.OnBossFloor {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	if player.BossKilled {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.BossKilled = true
	player.BossKillDuration = event.TimeSeconds - player.BossEnterTime

	return formatKilledBoss(event.TimeSeconds, player.ID)
}

func (p *Processor) handleLeaveDungeon(player *domain.Player, event domain.Event) string {
	if !player.Started {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.Finished = true
	player.FinishedAt = event.TimeSeconds

	if player.ClearedFloors == p.normalFloors() && player.BossKilled {
		player.State = domain.StateSuccess
	} else {
		player.State = domain.StateFail
	}

	return formatLeftDungeon(event.TimeSeconds, player.ID)
}

func (p *Processor) handleCannotContinue(player *domain.Player, event domain.Event) string {
	player.Finished = true
	player.FinishedAt = event.TimeSeconds
	player.State = domain.StateDisqual

	return formatCannotContinue(event.TimeSeconds, player.ID, event.Extra)
}

func (p *Processor) handleRestoreHealth(player *domain.Player, event domain.Event) string {
	if !p.isActive(player) {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	value, err := strconv.Atoi(event.Extra)
	if err != nil {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.Health += value
	if player.Health > 100 {
		player.Health = 100
	}

	return formatRestoredHealth(event.TimeSeconds, player.ID, value)
}

func (p *Processor) handleReceiveDamage(player *domain.Player, event domain.Event) []string {
	if !p.isActive(player) {
		return []string{
			formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID),
		}
	}

	value, err := strconv.Atoi(event.Extra)
	if err != nil {
		return []string{
			formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID),
		}
	}

	player.Health -= value

	output := []string{
		formatReceivedDamage(event.TimeSeconds, player.ID, value),
	}

	if player.Health <= 0 {
		player.Health = 0
		player.Finished = true
		player.FinishedAt = event.TimeSeconds
		player.State = domain.StateFail

		output = append(output, formatDead(event.TimeSeconds, player.ID))
	}

	return output
}
