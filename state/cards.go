package state

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

func (ths Gamestate) Draw(p constant.Player, n int) (Gamestate, []card.Card) {
	state := ths.Players[p]
	cards := state.Deck.PopX(n)
	state.Hand = append(state.Hand, cards...)
	ths.Players[p] = state
	return ths, cards
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

func (ths Gamestate) AddEnergy(p constant.Player, mid string, eid string) Gamestate {
	if ths.HasAttachedEnergy {
		panic("Already attached energy this turn")
	}
	state := ths.Players[p]

	var energy card.Card
	if c, ok := state.Hand.Remove(eid); ok {
		if c.CardType() != card.Energy {
			panic("Not energy card!")
		}
		energy = c
	} else {
		panic("Attempted to play card not in hand")
	}

	monster := state.getMonsterPointer(mid)
	if monster == nil {
		panic("Attempted to play card on nonexistant monster")
	}
	monster.Energy.Add(energy)

	ths.Players[p] = state
	ths.HasAttachedEnergy = true
	return ths
}
