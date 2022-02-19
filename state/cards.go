package state

import (
	"errors"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

func (ths Gamestate) Draw(p constant.Player, n int) (Gamestate, []Event) {
	state := ths.Players[p]
	cards := state.Deck.PopX(n)
	state.Hand = append(state.Hand, cards...)
	ths.Players[p] = state

	events := []Event{}
	for _, card := range cards {
		events = append(events, ECardDraw{card})
	}
	return ths, events
}

func (ths Gamestate) Shuffle(p constant.Player) Gamestate {
	ths.Players[p].Deck.Shuffle()
	return ths
}

func (ths Gamestate) PlaceHandOnDeck(p constant.Player) Gamestate {
	state := ths.Players[p]
	state.Deck = append(state.Deck, state.Hand...)
	state.Hand = card.CardGroup{}
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) AddEnergy(p constant.Player, mid string, eid string) (Gamestate, []Event, error) {
	if ths.HasAttachedEnergy {
		return ths, []Event{}, errors.New("already attached energy this turn")
	}
	state := ths.Players[p]

	var energy card.Card
	if c, ok := state.Hand.Remove(eid); ok {
		if c.CardType() != card.Energy {
			return ths, []Event{}, errors.New("not energy card")
		}
		energy = c
	} else {
		return ths, []Event{}, errors.New("attempted to play card not in hand")
	}

	monster := state.getMonsterPointer(mid)
	if monster == nil {
		return ths, []Event{}, errors.New("attempted to play card on nonexistant monster")
	}
	monster.Energy.Add(energy)

	ths.Players[p] = state
	ths.HasAttachedEnergy = true
	return ths, []Event{EAttachEnergy{Player: p, Energy: energy, Target: monster.Card}}, nil
}
