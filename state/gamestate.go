package state

import (
	"fmt"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

type Gamestate struct {
	Players map[constant.Player]Playerstate
	Turn    constant.Player
	Winner  *constant.Player
	IsDealt bool

	HasAttachedEnergy bool
}

type Playerstate struct {
	Deck    card.CardStack
	Hand    card.CardGroup
	Discard card.CardStack
	Prizes  card.CardStack

	Active         Cardstate
	HasInitialized bool
	HasActive      bool
	Bench          []Cardstate
}

type Cardstate struct {
	Card    card.Card
	Energy  card.CardGroup
	Damage  int
	Protect bool
}

func (ths Cardstate) String() string {
	hp := ths.Card.(card.MonsterCard).HP
	return fmt.Sprintf("%s HP %d/%d %s", ths.Card.String(), hp-ths.Damage, hp, ths.Energy.String())
}

func (ths Gamestate) Attack(p constant.Player, attackIndex int) Gamestate {
	t := constant.OtherPlayer(p)
	attack := ths.Players[p].Active.Card.(card.MonsterCard).Attacks[attackIndex]

	state := ths.Players[t]
	if !state.Active.Protect {
		state.Active.Damage += attack.Damage
	}
	ths.Players[t] = state

	if attack.Effect != nil {
		ths = LoadEffect(attack.Effect.Id, attack.Effect.Parameters).Apply(t, ths)
	}

	return ths.TurnTransition(p)
}

func (ths Gamestate) TurnTransition(p constant.Player) Gamestate {
	ths.Turn = constant.OtherPlayer(p)
	ths.HasAttachedEnergy = false

	state := ths.Players[ths.Turn]
	state.Active.Protect = false
	ths.Players[ths.Turn] = state

	return ths
}

func (ths Gamestate) SwitchTo(p constant.Player, benchIndex int) Gamestate {
	state := ths.Players[p]
	old := state.Active
	state.Active = state.Bench[benchIndex]
	state.Bench[benchIndex] = old
	ths.Players[p] = state
	return ths
}

func (ths Gamestate) SwitchDead(p constant.Player, benchIndex int) Gamestate {
	state := ths.Players[p]
	state.Active = state.Bench[benchIndex]
	state.Bench[benchIndex] = state.Bench[len(state.Bench)-1]
	state.Bench = state.Bench[:len(state.Bench)-1]
	ths.Players[p] = state

	ostate := ths.Players[constant.OtherPlayer(p)]
	ostate.Hand = append(ostate.Hand, ostate.Prizes.PopX(1)...)
	ths.Players[constant.OtherPlayer(p)] = ostate
	return ths
}

// PLACEHOLDER
func (ths Gamestate) Display() {
	fmt.Printf("#\t")
	for _, c := range ths.Players[constant.Player2].Bench {
		fmt.Printf("%s\t", c.Card.String())
	}
	fmt.Print("\n")
	fmt.Printf("#\n#\t%s\n#\n", ths.Players[constant.Player2].Active.String())
	fmt.Println("######################################################")
	fmt.Printf("#\n#\t%s\n#\n", ths.Players[constant.Player1].Active.String())
	fmt.Printf("#\t")
	for _, c := range ths.Players[constant.Player1].Bench {
		fmt.Printf("%s\t", c.Card.String())
	}
	fmt.Print("\n")
}
