package machine

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/state"
	"github.com/Danice123/cardmon/utils"
)

type Handler interface {
	Handle(Event)
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

func (ths *GameMachine) Event(event Event) {
	ths.handlers[constant.Player1].Handle(event)
	ths.handlers[constant.Player2].Handle(event)
}

func (ths *GameMachine) Start(game state.Gamestate) {
	ths.Current = game
	ths.Current = ths.Current.DealNewGame()
	ths.playInitialCards(constant.Player1)
	ths.playInitialCards(constant.Player2)

	e := Coinflip{Message: "if heads player 1 goes first"}
	if utils.Coinflip() {
		e.Outcome = "heads"
		ths.Current.Turn = constant.Player1
	} else {
		e.Outcome = "tails"
		ths.Current.Turn = constant.Player2
	}
	ths.Event(e)

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
	for {
		choice, _ := ths.handlers[player].AskCardFromHand("Choose a Basic monster to be your active.", false)
		if c, ok := ths.Current.Players[player].Hand.Get(choice); ok {
			if c.CardType() == card.Monster && c.(card.MonsterCard).Stage == 1 {
				ths.Current = ths.Current.PlayBasicFromHand(player, choice)
				break
			}
			ths.handlers[player].Alert("Card chosen was not a basic monster card")
		} else {
			ths.handlers[player].Alert("Card chosen was not in hand")
		}
	}

	for {
		choice, ok := ths.handlers[player].AskCardFromHand("Choose any Basic monsters to be placed on your bench.", true)
		if !ok {
			break
		}
		if c, ok := ths.Current.Players[player].Hand.Get(choice); ok {
			if c.CardType() == card.Monster && c.(card.MonsterCard).Stage == 1 {
				ths.Current = ths.Current.PlayBasicFromHand(player, choice)
			} else {
				ths.handlers[player].Alert("Card chosen was not a basic monster card")
			}
		} else {
			ths.handlers[player].Alert("Card chosen was not in hand")
		}
	}
}

func (ths *GameMachine) turn(p constant.Player) {
	var drawn []card.Card
	ths.Current, drawn = ths.Current.Draw(p, 1)
	for _, card := range drawn {
		ths.handlers[p].Handle(CardDraw{card})
	}

turnLoop:
	for ths.Current.Turn == p {
		action, choice := ths.handlers[p].AskForAction()
		switch action {
		case "Hand":
			if c, ok := ths.Current.Players[p].Hand.Get(choice); ok {
				ths.playCardFromHand(p, c)
			} else {
				ths.handlers[p].Alert("Card chosen was not in hand")
			}
		case "Attack":
			if attack, ok := ths.Current.Players[p].Active.Card.(card.MonsterCard).GetAttack(choice); ok {
				if attack.CheckCost(ths.Current.Players[p].Active.Energy) {
					ths.Current = ths.Current.Attack(p, choice)
					if ths.Current.Winner != nil {
						break turnLoop
					}

					opp := constant.OtherPlayer(p)
					if !ths.Current.Players[opp].HasActive {
						c, _ := ths.handlers[opp].AskTargetMonster("Choose monster to replace dead one.", false, false)
						ths.Current = ths.Current.SwitchDead(opp, c)
					}

					ths.Current = ths.Current.TurnTransition(p)
				} else {
					ths.handlers[p].Alert("Insufficient energy")
				}
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

func (ths *GameMachine) playCardFromHand(p constant.Player, c card.Card) {
	switch c.CardType() {
	case card.Energy:
		if ths.Current.HasAttachedEnergy {
			ths.handlers[p].Alert("Already attached energy this turn")
			break
		}
		if target, ok := ths.handlers[p].AskTargetMonster("Which monster to attach energy to?", true, true); ok {
			benchIndex, benched := ths.Current.Players[p].GetBenchIndex(target)
			if ths.Current.Players[p].Active.Card.Id() == target || benched {
				ths.Current = ths.Current.AddEnergy(p, target, c.Id())

				if benched {
					ths.Event(AttachEnergy{Player: p, Energy: c, Target: ths.Current.Players[p].Bench[benchIndex].Card})
				} else {
					ths.Event(AttachEnergy{Player: p, Energy: c, Target: ths.Current.Players[p].Active.Card})
				}
			} else {
				ths.handlers[p].Alert("Target monster doesn't exist")
			}
		}
	case card.Monster:
		if c.(card.MonsterCard).Stage == 1 {
			ths.Current = ths.Current.PlayBasicFromHand(p, c.Id())
		}
		// else {
		// handle evolutions
		//}
	case card.Trainer:
	}
}
