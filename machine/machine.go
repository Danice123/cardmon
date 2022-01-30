package machine

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/state"
	"github.com/Danice123/cardmon/utils"
)

type Handler interface {
	Alert(string)
	AskCardFromHand(string, bool) (int, bool)
	AskTargetMonster(string, bool) (int, bool)
	AskTargetBenchMonster(string, bool) (int, bool)
	AskForAction() (string, int)
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

func (ths *GameMachine) Start(game state.Gamestate) {
	ths.Current = game
	ths.Current = ths.Current.DealNewGame()
	ths.playInitialCards(constant.Player1)
	ths.playInitialCards(constant.Player2)

	ths.AlertBoth("Coin flip to see who goes first")
	if utils.Coinflip() {
		ths.Current.Turn = constant.Player2
	} else {
		ths.Current.Turn = constant.Player1
	}

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

func (ths *GameMachine) turn(p constant.Player) {
	ths.Current = ths.Current.Draw(p, 1)

turnLoop:
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

				opp := constant.OtherPlayer(p)
				if ths.Current.Players[opp].Active.Damage >= ths.Current.Players[opp].Active.Card.(card.MonsterCard).HP {
					if len(ths.Current.Players[opp].Bench) == 0 {
						ths.handlers[opp].Alert("You have run out of monsters on the field!")
						ths.Current.Winner = &p
						break turnLoop
					}
					c, _ := ths.handlers[opp].AskTargetBenchMonster("Choose monster to replace dead one.", false)
					ths.Current = ths.Current.SwitchDead(opp, c)
					if len(ths.Current.Players[p].Prizes) == 0 {
						ths.handlers[p].Alert("You have drawn all your prizes!")
						ths.Current.Winner = &p
						break turnLoop
					}
				}
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
