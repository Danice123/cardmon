package state

import (
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/utils"
)

type Effect interface {
	Apply(constant.Player, Gamestate) (Gamestate, []Event)
}

type effect struct {
	Id         string
	Parameters map[string]interface{}
}

type Coinflip struct {
	Description string
	Heads       *effect
	Tails       *effect
}

func (ths Coinflip) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	var outcome string
	events := []Event{}
	if utils.Coinflip() {
		outcome = "Heads"
		if ths.Heads != nil {
			gs, events = LoadEffect(ths.Heads.Id, ths.Heads.Parameters).Apply(t, gs)
		}
	} else {
		outcome = "Tails"
		if ths.Tails != nil {
			gs, events = LoadEffect(ths.Tails.Id, ths.Tails.Parameters).Apply(t, gs)
		}
	}
	return gs, append([]Event{ECoinflip{Message: ths.Description, Outcome: outcome}}, events...)
}

type Protect struct{}

func (ths Protect) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	state := gs.Players[constant.OtherPlayer(t)]
	state.Active.Protect = true
	gs.Players[constant.OtherPlayer(t)] = state
	return gs, []Event{}
}

type Selfdamage struct {
	Damage int
}

func (ths Selfdamage) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	state := gs.Players[constant.OtherPlayer(t)]
	state.Active.Damage += ths.Damage
	gs.Players[constant.OtherPlayer(t)] = state
	gs = gs.checkStateOfActive(constant.OtherPlayer(t))
	return gs, []Event{EDamage{Monster: state.Active.Card, Amount: ths.Damage}}
}
