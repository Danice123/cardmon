package machine

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/state"
)

type Handler interface {
	Handle(state.Event)
	Alert(string)
	AskCardFromHand(message string, cancelable bool) (string, bool)
	AskTargetMonster(message string, showActive bool, cancelable bool) (string, bool)
	AskForAction() (string, string)
}

type GameMachine struct {
	handlers map[constant.Player]Handler

	Current state.Gamestate
}

func (ths *GameMachine) RegisterHandler(player constant.Player, handler Handler) {
	if ths.handlers == nil {
		ths.handlers = map[constant.Player]Handler{}
	}
	ths.handlers[player] = handler
}

func (ths *GameMachine) AlertBoth(message string) {
	ths.handlers[constant.Player1].Alert(message)
	ths.handlers[constant.Player2].Alert(message)
}

func (ths *GameMachine) Events(events []state.Event) {
	for _, e := range events {
		ths.handlers[constant.Player1].Handle(e)
		ths.handlers[constant.Player2].Handle(e)
	}
}

func (ths *GameMachine) Start(game state.Gamestate) {
	ths.Current = game
	ths.Current = ths.Current.DealNewGame()
	ths.playInitialCards(constant.Player1)
	ths.playInitialCards(constant.Player2)

	var events []state.Event
	var err error
	ths.Current, events, err = ths.Current.StartGame()
	if err != nil {
		panic(err)
	}
	ths.Events(events)

	for {
		if len(ths.Current.Players[ths.Current.Turn].Deck) == 0 {
			ths.handlers[ths.Current.Turn].Alert("You have run out of cards to draw!")
			opp := constant.OtherPlayer(ths.Current.Turn)
			ths.Current.Winner = &opp
		} else {
			ths.turn(ths.Current.Turn)
		}

		if ths.Current.Winner != nil {
			ths.handlers[*ths.Current.Winner].Alert("You win")
			ths.handlers[constant.OtherPlayer(*ths.Current.Winner)].Alert("You lose")
			break
		}
	}
}

func (ths *GameMachine) playInitialCards(player constant.Player) {
	var err error
	for {
		choice, _ := ths.handlers[player].AskCardFromHand("Choose a Basic monster to be your active.", false)
		ths.Current, _, err = ths.Current.PlayBasicFromHand(player, choice)
		if err != nil {
			ths.handlers[player].Alert(err.Error())
		} else {
			break
		}
	}

	for {
		choice, ok := ths.handlers[player].AskCardFromHand("Choose any Basic monsters to be placed on your bench.", true)
		if !ok {
			break
		}
		ths.Current, _, err = ths.Current.PlayBasicFromHand(player, choice)
		if err != nil {
			ths.handlers[player].Alert(err.Error())
		}
	}
}

func (ths *GameMachine) turn(p constant.Player) {
	var err error
	var events []state.Event
	ths.Current, events = ths.Current.Draw(p, 1)
	ths.handlers[p].Handle(events[0])

turnLoop:
	for ths.Current.Turn == p {
		action, choice := ths.handlers[p].AskForAction()
		switch action {
		case "Hand":
			if c, ok := ths.Current.Players[p].Hand.Get(choice); ok {
				switch c.CardType() {
				case card.Energy:
					if target, ok := ths.handlers[p].AskTargetMonster("Which monster to attach energy to?", true, true); ok {
						ths.Current, events, err = ths.Current.AddEnergy(p, target, c.Id())
						if err != nil {
							ths.handlers[p].Alert(err.Error())
							continue turnLoop
						}
						ths.Events(events)
					}
				case card.Monster:
					if c.(card.MonsterCard).Stage == 1 {
						ths.Current, events, err = ths.Current.PlayBasicFromHand(p, c.Id())
						if err != nil {
							ths.handlers[p].Alert(err.Error())
							continue turnLoop
						}
						ths.Events(events)
					}
					// else {
					// handle evolutions
					//}
				case card.Trainer:
				}
			} else {
				ths.handlers[p].Alert("Card chosen was not in hand")
			}
		case "Attack":
			ths.Current, events, err = ths.Current.Attack(p, choice)
			if err != nil {
				ths.handlers[p].Alert(err.Error())
				continue turnLoop
			}
			ths.Events(events)
			if ths.Current.Winner != nil {
				break turnLoop
			}

			opp := constant.OtherPlayer(p)
			if !ths.Current.Players[opp].HasActive {
				c, _ := ths.handlers[opp].AskTargetMonster("Choose monster to replace dead one.", false, false)
				ths.Current = ths.Current.SwitchDead(opp, c)
			}

			ths.Current = ths.Current.TurnTransition(p)
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
