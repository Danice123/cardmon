package machine

import (
	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
)

type Event interface {
	Type() EventType
}

type Coinflip struct {
	Message string
	Outcome string
}

type CardDraw struct {
	Card card.Card
}

type AttachEnergy struct {
	Energy card.Card
	Target card.Card
	Player constant.Player
}
