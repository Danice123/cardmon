package state

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

type Event interface {
	Type() EventType
}

type ECoinflip struct {
	Message string
	Outcome string
}

type ECardDraw struct {
	Card card.Card
}

type EAttachEnergy struct {
	Energy card.Card
	Target card.Card
	Player constant.Player
}

type EAddToBench struct {
	Player  constant.Player
	Monster card.Card
}

type EDamage struct {
	Monster card.Card
	Amount  int
}
