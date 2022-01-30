package state

import (
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/utils"
)

type Effect interface {
	Apply(constant.Player, Gamestate) Gamestate
}

type effect struct {
	Id         string
	Parameters map[string]interface{}
}

type Coinflip struct {
	Heads *effect
	Tails *effect
}

func (ths Coinflip) Apply(t constant.Player, gs Gamestate) Gamestate {
	if utils.Coinflip() {
		if ths.Heads != nil {
			return LoadEffect(ths.Heads.Id, ths.Heads.Parameters).Apply(t, gs)
		}
	} else {
		if ths.Tails != nil {
			return LoadEffect(ths.Tails.Id, ths.Tails.Parameters).Apply(t, gs)
		}
	}
	return gs
}

type Protect struct{}

func (ths Protect) Apply(t constant.Player, gs Gamestate) Gamestate {
	state := gs.Players[constant.OtherPlayer(t)]
	state.Active.Protect = true
	gs.Players[constant.OtherPlayer(t)] = state
	return gs
}
