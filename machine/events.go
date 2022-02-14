package machine

import "github.com/Danice123/cardmon/card"

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
	Energy card.EnergyCard
	Target int
}
