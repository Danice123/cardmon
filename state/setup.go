package state

import (
	"errors"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/utils"
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
	if ths.IsDealt {
		panic("Game is already running")
	}
	ths = ths.Shuffle(constant.Player1)
	ths = ths.Shuffle(constant.Player2)
	ths, _ = ths.Draw(constant.Player1, 7)
	ths, _ = ths.Draw(constant.Player2, 7)
	ths = ths.checkForBasic(constant.Player1)
	ths = ths.checkForBasic(constant.Player2)
	ths = ths.placePrizes(constant.Player1)
	ths = ths.placePrizes(constant.Player2)
	ths.IsDealt = true
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
			ths, _ = ths.Draw(p, 7)
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
				Card: c,
			}
			state.HasActive = true
		} else {
			state.Bench = append(state.Bench, Cardstate{
				Card: c,
			})
			events = append(events, EAddToBench{Player: p, Monster: c})
		}
		ths.Players[p] = state
		return ths, events, nil
	} else {
		return ths, events, errors.New("attempted to play card not in hand")
	}
}

func (ths Gamestate) StartGame() (Gamestate, []Event, error) {
	for p, state := range ths.Players {
		if state.HasInitialized {
			return ths, []Event{}, errors.New("player has already been initialized")
		}
		if state.HasActive {
			state.HasInitialized = true
		} else {
			return ths, []Event{}, errors.New("player has not placed active monster")
		}
		ths.Players[p] = state
	}

	e := ECoinflip{Message: "if heads player 1 goes first"}
	if utils.Coinflip() {
		e.Outcome = "heads"
		ths.Turn = constant.Player1
	} else {
		e.Outcome = "tails"
		ths.Turn = constant.Player2
	}

	return ths, []Event{e}, nil
}
