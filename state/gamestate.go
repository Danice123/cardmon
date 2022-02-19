package state

import (
	"errors"
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
	Card     card.Card
	Energy   card.CardGroup
	Damage   int
	Statuses StatusList
	Protect  bool
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

func (ths Gamestate) Attack(p constant.Player, aid string) (Gamestate, []Event, error) {
	for _, status := range ths.Players[p].Active.Statuses {
		if resp, ok := status.CanMonsterAttack(); !ok {
			return ths, []Event{}, errors.New(resp)
		}
	}

	events := []Event{}
	t := constant.OtherPlayer(p)
	if attack, ok := ths.Players[p].Active.Card.(card.MonsterCard).GetAttack(aid); ok {
		if attack.CheckCost(ths.Players[p].Active.Energy) {
			for _, effect := range attack.Effects {
				var effectEvents []Event
				ths, effectEvents = LoadEffect(effect.Id, effect.Parameters).Apply(t, ths)
				events = append(events, effectEvents...)
			}
			ths = ths.checkStateOfActive(t)
			return ths, events, nil
		} else {
			return ths, events, errors.New("insufficient energy")
		}

	} else {
		return ths, events, errors.New("invalid attack id")
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

func (ths Gamestate) PlayBasicFromHand(p constant.Player, cid string) (Gamestate, []Event, error) {
	events := []Event{}
	state := ths.Players[p]
	if c, ok := state.Hand.Remove(cid); ok {
		if c.CardType() != card.Monster {
			return ths, events, errors.New("attempted to play non-monster card to field")
		}
		if c.(card.MonsterCard).Stage != 1 {
			return ths, events, errors.New("attempted to play non-basic monster card to field")
		}
		if !state.HasActive && !state.HasInitialized {
			state.Active = Cardstate{
				Card:     c,
				Statuses: StatusList{},
			}
			state.HasActive = true
		} else {
			state.Bench = append(state.Bench, Cardstate{
				Card:     c,
				Statuses: StatusList{},
			})
			events = append(events, EAddToBench{Player: p, Monster: c})
		}
		ths.Players[p] = state
		return ths, events, nil
	} else {
		return ths, events, errors.New("attempted to play card not in hand")
	}
}

func (ths Gamestate) TurnTransition(p constant.Player) (Gamestate, []Event) {
	ths.Turn = constant.OtherPlayer(p)
	ths.HasAttachedEnergy = false

	events := []Event{}
	for target, state := range ths.Players {
		for _, status := range state.Active.Statuses {
			var e []Event
			ths, e = status.OnTurnChange(ths, target)
			events = append(events, e...)
		}
	}

	state := ths.Players[ths.Turn]
	state.Active.Protect = false
	ths.Players[ths.Turn] = state

	return ths, events
}

func (ths Gamestate) SwitchTo(p constant.Player, mid string) (Gamestate, []Event, error) {
	for _, status := range ths.Players[p].Active.Statuses {
		if resp, ok := status.CanMonsterAttack(); !ok {
			return ths, []Event{}, errors.New(resp)
		}
	}

	state := ths.Players[p]
	if benchIndex, ok := state.GetBenchIndex(mid); ok {
		if ths.Players[p].Active.Card.(card.MonsterCard).Retreat <= len(ths.Players[p].Active.Energy) {
			old := state.Active
			state.Active = state.Bench[benchIndex]
			state.Bench[benchIndex] = old
			ths.Players[p] = state
			return ths, []Event{}, nil
		} else {
			return ths, []Event{}, errors.New("insufficient energy")
		}
	} else {
		return ths, []Event{}, errors.New("attempted to bring out nonexistant monster")
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
