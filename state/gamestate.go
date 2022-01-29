package state

import (
	"fmt"

	"github.com/Danice123/cardmon/card"
)

type Player string

const Player1 = Player("one")
const Player2 = Player("two")

func OtherPlayer(p Player) Player {
	if p == Player1 {
		return Player2
	} else {
		return Player1
	}
}

type Gamestate struct {
	Players           map[Player]Playerstate
	Turn              Player
	Winner            *Player
	HasAttachedEnergy bool
}

type Playerstate struct {
	Deck    card.CardStack
	Hand    card.CardGroup
	Discard card.CardStack
	Prizes  card.CardStack

	Active    Cardstate
	HasActive bool
	Bench     []Cardstate
}

type Cardstate struct {
	Card   card.Card
	Energy card.CardGroup
	Damage int
}

func (ths Cardstate) String() string {
	hp := ths.Card.(card.MonsterCard).HP
	return fmt.Sprintf("%s HP %d/%d %s", ths.Card.String(), hp-ths.Damage, hp, ths.Energy.String())
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

func (ths Gamestate) PlayBasicFromHand(p Player, handIndex int) Gamestate {
	state := ths.Players[p]
	c := state.Hand.Remove(handIndex)
	if c.CardType() != card.Monster {
		panic("Attempted to play non-monster card to field")
	}
	if c.(card.MonsterCard).Stage != 1 {
		panic("Attempted to play non-basic monster card to field")
	}
	if !state.HasActive {
		state.Active = Cardstate{
			Card: c,
		}
		state.HasActive = true
	} else {
		state.Bench = append(state.Bench, Cardstate{
			Card: c,
		})
	}
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) PlacePrizes(p Player) Gamestate {
	state := ths.Players[p]
	state.Prizes = state.Deck.PopX(6)
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) Draw(p Player, n int) Gamestate {
	state := ths.Players[p]
	state.Hand = append(state.Hand, state.Deck.PopX(n)...)
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) Shuffle(p Player) Gamestate {
	ths.Players[p].Deck.Shuffle()
	return ths
}

func (ths Gamestate) PlaceHandOnDeck(p Player) Gamestate {
	state := ths.Players[p]
	state.Deck = append(state.Deck, state.Hand...)
	state.Hand = card.CardGroup{}
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) AddEnergyToActive(p Player, handIndex int) Gamestate {
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

func (ths Gamestate) AddEnergyToBench(p Player, handIndex int, benchIndex int) Gamestate {
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

func (ths Gamestate) Attack(p Player, attackIndex int) Gamestate {
	state := ths.Players[OtherPlayer(p)]
	state.Active.Damage += ths.Players[p].Active.Card.(card.MonsterCard).Attacks[attackIndex].Damage
	ths.Players[OtherPlayer(p)] = state
	return ths.TurnTransition(p)
}

func (ths Gamestate) SwitchTo(p Player, benchIndex int) Gamestate {
	state := ths.Players[p]
	old := state.Active
	state.Active = state.Bench[benchIndex]
	state.Bench[benchIndex] = old
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) SwitchDead(p Player, benchIndex int) Gamestate {
	state := ths.Players[p]
	state.Active = state.Bench[benchIndex]
	state.Bench[benchIndex] = state.Bench[len(state.Bench)-1]
	state.Bench = state.Bench[:len(state.Bench)-1]
	ths.Players[p] = state

	ostate := ths.Players[OtherPlayer(p)]
	ostate.Hand = append(ostate.Hand, ostate.Prizes.PopX(1)...)
	ths.Players[OtherPlayer(p)] = ostate
	return ths
}

func (ths Gamestate) TurnTransition(p Player) Gamestate {
	ths.Turn = OtherPlayer(p)
	ths.HasAttachedEnergy = false
	return ths
}

// PLACEHOLDER
func (ths Gamestate) Display() {
	fmt.Printf("#\t")
	for _, c := range ths.Players[Player2].Bench {
		fmt.Printf("%s\t", c.Card.String())
	}
	fmt.Print("\n")
	fmt.Printf("#\n#\t%s\n#\n", ths.Players[Player2].Active.String())
	fmt.Println("######################################################")
	fmt.Printf("#\n#\t%s\n#\n", ths.Players[Player1].Active.String())
	fmt.Printf("#\t")
	for _, c := range ths.Players[Player1].Bench {
		fmt.Printf("%s\t", c.Card.String())
	}
	fmt.Print("\n")
}
