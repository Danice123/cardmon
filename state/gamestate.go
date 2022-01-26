package state

import (
	"fmt"

	"github.com/Danice123/cardmon/card"
)

type Player string

const Player1 = Player("one")
const Player2 = Player("two")

type Gamestate struct {
	Players map[Player]Playerstate
	Turn    Player
}

type Playerstate struct {
	Deck    card.CardStack
	Hand    card.CardGroup
	Discard card.CardStack

	Active    Cardstate
	HasActive bool
	Bench     []Cardstate
}

type Cardstate struct {
	Card card.Card
}

func NewGame(deck1 card.CardStack, deck2 card.CardStack) Gamestate {
	return Gamestate{
		Players: map[Player]Playerstate{
			Player1: {
				Deck: deck1,
			},
			Player2: {
				Deck: deck2,
			},
		},
	}
}

func (ths Gamestate) Shuffle(p Player) Gamestate {
	ths.Players[p].Deck.Shuffle()
	return ths
}

func (ths Gamestate) PlayBasicFromHand(player Player, handIndex int) Gamestate {
	state := ths.Players[player]
	card := state.Hand.Remove(handIndex)
	if !state.HasActive {
		state.Active = Cardstate{
			Card: card,
		}
	} else {
		state.Bench = append(state.Bench, Cardstate{
			Card: card,
		})
	}
	ths.Players[player] = state
	return ths
}

func (ths Gamestate) Deal() Gamestate {
	for player, state := range ths.Players {
		state.Hand = state.Deck.PopX(7)
		ths.Players[player] = state
	}
	return ths
}

func (ths Gamestate) Display() {
	fmt.Printf("#\n#\t%s\n#\n", ths.Players[Player2].Active.Card.ShortText())
	fmt.Println("######################################################")
	fmt.Printf("#\n#\t%s\n#\n", ths.Players[Player1].Active.Card.ShortText())
}
