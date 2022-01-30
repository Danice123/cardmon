package state

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

func NewGame(deck1 card.CardStack, deck2 card.CardStack) Gamestate {
	return Gamestate{
		Players: map[constant.Player]Playerstate{
			constant.Player1: {
				Deck: deck1,
			},
			constant.Player2: {
				Deck: deck2,
			},
		},
	}
}

func (ths Gamestate) DealNewGame() Gamestate {
	if !ths.IsDealt {
		panic("Game is already running")
	}
	ths = ths.Shuffle(constant.Player1)
	ths = ths.Shuffle(constant.Player2)
	ths = ths.Draw(constant.Player1, 7)
	ths = ths.Draw(constant.Player2, 7)
	ths = ths.checkForBasic(constant.Player1)
	ths = ths.checkForBasic(constant.Player2)
	ths = ths.placePrizes(constant.Player1)
	ths = ths.placePrizes(constant.Player2)
	ths.IsDealt = false
	return ths
}

func (ths Gamestate) checkForBasic(p constant.Player) Gamestate {
	for {
		var hasBasic bool
		for _, c := range ths.Players[p].Hand {
			if c.CardType() == card.Monster && c.(card.MonsterCard).Stage == 1 {
				hasBasic = true
				break
			}
		}
		if !hasBasic {
			ths = ths.PlaceHandOnDeck(p)
			ths = ths.Shuffle(p)
			ths = ths.Draw(p, 7)
			continue
		}
		break
	}
	return ths
}

func (ths Gamestate) placePrizes(p constant.Player) Gamestate {
	state := ths.Players[p]
	state.Prizes = state.Deck.PopX(6)
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) PlayBasicFromHand(p constant.Player, handIndex int) Gamestate {
	state := ths.Players[p]
	c := state.Hand.Remove(handIndex)
	if c.CardType() != card.Monster {
		panic("Attempted to play non-monster card to field")
	}
	if c.(card.MonsterCard).Stage != 1 {
		panic("Attempted to play non-basic monster card to field")
	}
	if !state.HasActive && !state.HasInitialized {
		state.Active = Cardstate{
			Card: c,
		}
		state.HasActive = true
		state.HasInitialized = true
	} else {
		state.Bench = append(state.Bench, Cardstate{
			Card: c,
		})
	}
	ths.Players[p] = state
	return ths
}
