package state

import (
	"fmt"
	"strings"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

type Gamestate struct {
	Players map[constant.Player]Playerstate
	Turn    constant.Player
	Winner  *constant.Player
	IsDealt bool

	HasAttachedEnergy bool
}

type Playerstate struct {
	Deck    card.CardStack
	Hand    card.CardGroup
	Discard card.CardStack
	Prizes  card.CardStack

	Active         Cardstate
	HasInitialized bool
	HasActive      bool
	Bench          []Cardstate
}

func (ths *Playerstate) getMonsterPointer(mid string) *Cardstate {
	var monster *Cardstate
	if ths.Active.Card.Id() == mid {
		monster = &ths.Active
	} else {
		for i, m := range ths.Bench {
			if m.Card.Id() == mid {
				monster = &ths.Bench[i]
				break
			}
		}
	}
	return monster
}

func (ths Playerstate) GetBenchIndex(mid string) (int, bool) {
	index := -1
	for i, m := range ths.Bench {
		if m.Card.Id() == mid {
			index = i
			break
		}
	}
	return index, index != -1
}

type Cardstate struct {
	Card    card.Card
	Energy  card.CardGroup
	Damage  int
	Protect bool
}

func (ths Cardstate) String() string {
	hp := ths.Card.(card.MonsterCard).HP
	sb := strings.Builder{}
	for _, e := range ths.Energy {
		sb.WriteString(e.(card.EnergyCard).Pprint())
		sb.WriteString(" ")
	}
	return fmt.Sprintf("%s HP %d/%d %s", ths.Card.String(), hp-ths.Damage, hp, sb.String())
}

func (ths Gamestate) Attack(p constant.Player, aid string) Gamestate {
	t := constant.OtherPlayer(p)
	if attack, ok := ths.Players[p].Active.Card.(card.MonsterCard).GetAttack(aid); ok {
		state := ths.Players[t]
		if !state.Active.Protect {
			state.Active.Damage += attack.Damage
		}
		ths.Players[t] = state
		ths = ths.checkStateOfActive(t)

		if attack.Effect != nil {
			ths = LoadEffect(attack.Effect.Id, attack.Effect.Parameters).Apply(t, ths)
		}

		return ths
	} else {
		panic("Bad attack id")
	}
}

func (ths Gamestate) checkStateOfActive(p constant.Player) Gamestate {
	state := ths.Players[p]
	if state.Active.Damage >= state.Active.Card.(card.MonsterCard).HP {
		state.HasActive = false
		ths.Players[p] = state

		ostate := ths.Players[constant.OtherPlayer(p)]
		ostate.Hand = append(ostate.Hand, ostate.Prizes.PopX(1)...)
		ths.Players[constant.OtherPlayer(p)] = ostate

		if len(state.Bench) == 0 || len(ostate.Prizes) == 0 {
			win := constant.OtherPlayer(p)
			ths.Winner = &win
		}
	}
	return ths
}

func (ths Gamestate) TurnTransition(p constant.Player) Gamestate {
	ths.Turn = constant.OtherPlayer(p)
	ths.HasAttachedEnergy = false

	state := ths.Players[ths.Turn]
	state.Active.Protect = false
	ths.Players[ths.Turn] = state

	return ths
}

func (ths Gamestate) SwitchTo(p constant.Player, mid string) Gamestate {
	state := ths.Players[p]
	if benchIndex, ok := state.GetBenchIndex(mid); ok {
		old := state.Active
		state.Active = state.Bench[benchIndex]
		state.Bench[benchIndex] = old
		ths.Players[p] = state
		return ths
	} else {
		panic("Attempted to bring out nonexistant monster")
	}
}

func (ths Gamestate) SwitchDead(p constant.Player, mid string) Gamestate {
	state := ths.Players[p]
	if benchIndex, ok := state.GetBenchIndex(mid); ok {
		state.Active = state.Bench[benchIndex]
		state.Bench[benchIndex] = state.Bench[len(state.Bench)-1]
		state.Bench = state.Bench[:len(state.Bench)-1]
		state.HasActive = true
		ths.Players[p] = state
		return ths
	} else {
		panic("Attempted to bring out nonexistant monster")
	}
}
