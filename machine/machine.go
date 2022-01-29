package machine

import (
	"fmt"
	"math/rand"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/state"
)

func coinflip() bool {
	if rand.Intn(100) > 50 {
		return true
	} else {
		return false
	}
}

type Handler interface {
	Alert(string)
	AskCardFromHand(string, bool) (int, bool)
	AskTargetMonster(string, bool) (int, bool)
	AskForAction() (string, int)
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
	ths.Current = ths.Current.Draw(state.Player1, 7)
	ths.Current = ths.Current.Draw(state.Player2, 7)
	ths.checkForBasics(state.Player1)
	ths.checkForBasics(state.Player2)
	ths.playInitialCards(state.Player1)
	ths.playInitialCards(state.Player2)

	ths.AlertBoth("Placing six prizes")
	ths.Current = ths.Current.PlacePrizes(state.Player1)
	ths.Current = ths.Current.PlacePrizes(state.Player2)

	ths.AlertBoth("Coin flip to see who goes first")
	if coinflip() {
		ths.Current.Turn = state.Player2
	} else {
		ths.Current.Turn = state.Player1
	}

	for {
		// Check for game end by draw exhaust
		ths.turn(ths.Current.Turn)
		// Check for game end by prize or lack of active
	}
}

func (ths *GameMachine) checkForBasics(player state.Player) {
	for {
		var hasBasic bool
		for _, c := range ths.Current.Players[player].Hand {
			if c.CardType() == card.Monster && c.(card.MonsterCard).Stage == 1 {
				hasBasic = true
				break
			}
		}
		if !hasBasic {
			ths.AlertBoth(fmt.Sprintf("Player %s has no basic cards, reshuffling deck and drawing again", player))
			ths.Current = ths.Current.PlaceHandOnDeck(player)
			ths.Current = ths.Current.Shuffle(player)
			ths.Current = ths.Current.Draw(player, 7)
			continue
		}
		break
	}
}

func (ths *GameMachine) playInitialCards(player state.Player) {
	for {
		choice, _ := ths.handlers[player].AskCardFromHand("Choose a Basic monster to be your active.", false)
		c := ths.Current.Players[player].Hand[choice]
		if c.CardType() == card.Monster && c.(card.MonsterCard).Stage == 1 {
			ths.Current = ths.Current.PlayBasicFromHand(player, choice)
			break
		}
		ths.handlers[player].Alert("Card chosen was not a basic monster card")
	}

	for {
		choice, ok := ths.handlers[player].AskCardFromHand("Choose any Basic monsters to be placed on your bench.", true)
		if !ok {
			break
		}
		c := ths.Current.Players[player].Hand[choice]
		if c.CardType() == card.Monster && c.(card.MonsterCard).Stage == 1 {
			ths.Current = ths.Current.PlayBasicFromHand(player, choice)
		} else {
			ths.handlers[player].Alert("Card chosen was not a basic monster card")
		}
	}
}

func (ths *GameMachine) turn(p state.Player) {
	for ths.Current.Turn == p {
		ths.Current.Display()
		action, choice := ths.handlers[p].AskForAction()

		switch action {
		case "Hand":
			c := ths.Current.Players[p].Hand[choice]
			switch c.CardType() {
			case card.Energy:
				if ths.Current.HasAttachedEnergy {
					ths.handlers[p].Alert("Already attached energy this turn")
					break
				}
				if target, ok := ths.handlers[p].AskTargetMonster("Which monster to attach energy to?", true); ok {
					if target == 0 {
						ths.Current = ths.Current.AddEnergyToActive(p, choice)
					} else {
						ths.Current = ths.Current.AddEnergyToBench(p, choice, target-1)
					}
				}
			case card.Monster:
				if c.(card.MonsterCard).Stage == 1 {
					ths.Current = ths.Current.PlayBasicFromHand(p, choice)
				} // else {
				// handle evolutions
				//}
			case card.Trainer:
			}
		case "Attack":
			a := ths.Current.Players[p].Active.Card.(card.MonsterCard).Attacks[choice]
			if a.CheckCost(ths.Current.Players[p].Active.Energy) {
				ths.Current = ths.Current.Attack(p, choice)
			} else {
				ths.handlers[p].Alert("Insufficient energy")
			}
		case "Retreat":
			if ths.Current.Players[p].Active.Card.(card.MonsterCard).Retreat <= len(ths.Current.Players[p].Active.Energy) {
				ths.Current = ths.Current.SwitchTo(p, choice)
			} else {
				ths.handlers[p].Alert("Insufficient energy")
			}
		case "Pass":
			ths.Current = ths.Current.TurnTransition(p)
		}
	}
}
