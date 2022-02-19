package state

type EventType string

var COINFLIP = EventType("COINFLIP")

func (ths ECoinflip) Type() EventType {
	return COINFLIP
}

var CARDDRAW = EventType("CARDDRAW")

func (ths ECardDraw) Type() EventType {
	return CARDDRAW
}

var ATTACHENERGY = EventType("ATTACHENERGY")

func (ths EAttachEnergy) Type() EventType {
	return ATTACHENERGY
}

var ADDTOBENCH = EventType("ADDTOBENCH")

func (ths EAddToBench) Type() EventType {
	return ADDTOBENCH
}

var DAMAGE = EventType("DAMAGE")

func (ths EDamage) Type() EventType {
	return DAMAGE
}
