package state

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

func (ths Gamestate) Draw(p constant.Player, n int) Gamestate {
	state := ths.Players[p]
	state.Hand = append(state.Hand, state.Deck.PopX(n)...)
	ths.Players[p] = state
	return ths
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

func (ths Gamestate) AddEnergyToActive(p constant.Player, handIndex int) Gamestate {
	if ths.HasAttachedEnergy {
		panic("Already attached energy this turn")
	}
	state := ths.Players[p]
	c := state.Hand.Remove(handIndex)
	if c.CardType() != card.Energy {
		panic("Not energy card!")
	}
	state.Active.Energy.Add(c)
	ths.Players[p] = state
	ths.HasAttachedEnergy = true
	return ths
}

func (ths Gamestate) AddEnergyToBench(p constant.Player, handIndex int, benchIndex int) Gamestate {
	if ths.HasAttachedEnergy {
		panic("Already attached energy this turn")
	}
	state := ths.Players[p]
	c := state.Hand.Remove(handIndex)
	if c.CardType() != card.Energy {
		panic("Not energy card!")
	}
	state.Bench[benchIndex].Energy.Add(c)
	ths.Players[p] = state
	ths.HasAttachedEnergy = true
	return ths
}
