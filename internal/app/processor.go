package app

import (
	"github.com/shar1mo/dungeon-challenge/internal/config"
	"github.com/shar1mo/dungeon-challenge/internal/domain"
)

const (
	eventRegister       = 1
	eventEnterDungeon   = 2
	eventLeaveDungeon   = 8
	eventCannotContinue = 9
)

type Processor struct {
	cfg     config.Config
	players map[int]*domain.Player
}

func NewProcessor(cfg config.Config) (*Processor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Processor{
		cfg:     cfg,
		players: make(map[int]*domain.Player),
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
		case eventLeaveDungeon:
			output = append(output, p.handleLeaveDungeon(player, event))
		case eventCannotContinue:
			output = append(output, p.handleCannotContinue(player, event))
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

	return player
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

func (p *Processor) handleLeaveDungeon(player *domain.Player, event domain.Event) string {
	if !player.Started {
		return formatImpossibleMove(event.TimeSeconds, player.ID, event.EventID)
	}

	player.Finished = true
	player.FinishedAt = event.TimeSeconds
	player.State = domain.StateFail

	return formatLeftDungeon(event.TimeSeconds, player.ID)
}

func (p *Processor) handleCannotContinue(player *domain.Player, event domain.Event) string {
	player.Finished = true
	player.FinishedAt = event.TimeSeconds
	player.State = domain.StateDisqual

	return formatCannotContinue(event.TimeSeconds, player.ID, event.Extra)
}
