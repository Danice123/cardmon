package machine

import (
	"math/rand"

	"github.com/Danice123/cardmon/state"
)

type Handler interface {
	Alert(string)
	AskCardFromHand(string) int
}

type GameMachine struct {
	handlers map[state.Player]Handler

	Current state.Gamestate
}

func (ths *GameMachine) RegisterHandler(player state.Player, handler Handler) {
	if ths.handlers == nil {
		ths.handlers = map[state.Player]Handler{}
	}
	ths.handlers[player] = handler
}

func (ths *GameMachine) AlertBoth(message string) {
	ths.handlers[state.Player1].Alert(message)
	ths.handlers[state.Player2].Alert(message)
}

func (ths *GameMachine) Start(game state.Gamestate) {
	ths.Current = game
	ths.AlertBoth("Shuffling your deck")
	ths.Current = ths.Current.Shuffle(state.Player1)
	ths.Current = ths.Current.Shuffle(state.Player2)

	ths.AlertBoth("Each player will draw 7 cards")
	ths.Current = ths.Current.Deal()

	// Handle no basics

	card := ths.handlers[state.Player1].AskCardFromHand("Choose a Basic monster to be your active.")
	ths.Current = ths.Current.PlayBasicFromHand(state.Player1, card)
	card = ths.handlers[state.Player2].AskCardFromHand("Choose a Basic monster to be your active.")
	ths.Current = ths.Current.PlayBasicFromHand(state.Player2, card)

	ths.Current.Display()

	ths.AlertBoth("Coin flip to see who goes first")
	if rand.Intn(2) > 0 {
		ths.Current.Turn = state.Player2
	} else {
		ths.Current.Turn = state.Player1
	}
}
