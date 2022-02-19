package state

import (
	"github.com/Danice123/cardmon/card"
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

type Damage struct {
	Amount int
}

func (ths Damage) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	state := gs.Players[t]
	if !state.Active.Protect {
		damage := ths.Amount
		println(state.Active.Card.(card.MonsterCard).Id())
		println(gs.Players[constant.OtherPlayer(t)].Active.Card.(card.MonsterCard).Id())
		if state.Active.Card.(card.MonsterCard).Weakness == gs.Players[constant.OtherPlayer(t)].Active.Card.(card.MonsterCard).Type {
			damage += 30
		}
		state.Active.Damage += damage
		gs.Players[t] = state
		return gs, []Event{EDamage{Monster: state.Active.Card, Amount: damage}}
	}
	return gs, []Event{}
}

type Coinflip struct {
	Description string
	Heads       *effect
	Tails       *effect
}

func (ths Coinflip) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	outcome := utils.Coinflip()
	events := []Event{}
	if outcome {
		if ths.Heads != nil {
			gs, events = LoadEffect(ths.Heads.Id, ths.Heads.Parameters).Apply(t, gs)
		}
	} else {
		if ths.Tails != nil {
			gs, events = LoadEffect(ths.Tails.Id, ths.Tails.Parameters).Apply(t, gs)
		}
	}
	return gs, append([]Event{ECoinflip{Message: ths.Description, Outcome: outcome}}, events...)
}

type SleepEffect struct {
	Self bool
}

func (ths SleepEffect) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	target := t
	if ths.Self {
		target = constant.OtherPlayer(target)
	}

	state := gs.Players[target]
	state.Active.Statuses.Add(SleepStatus{})
	gs.Players[target] = state
	return gs, []Event{}
}

type Protect struct{}

func (ths Protect) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	state := gs.Players[constant.OtherPlayer(t)]
	state.Active.Protect = true
	gs.Players[constant.OtherPlayer(t)] = state
	return gs, []Event{}
}

type Selfdamage struct {
	Amount int
}

func (ths Selfdamage) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	state := gs.Players[constant.OtherPlayer(t)]
	state.Active.Damage += ths.Amount
	gs.Players[constant.OtherPlayer(t)] = state
	gs = gs.checkStateOfActive(constant.OtherPlayer(t))
	return gs, []Event{EDamage{Monster: state.Active.Card, Amount: ths.Amount}}
}

type Metronome struct{}

func (ths Metronome) Apply(t constant.Player, gs Gamestate) (Gamestate, []Event) {
	return gs, []Event{}
}
