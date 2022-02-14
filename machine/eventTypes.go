package machine

type EventType string

var COINFLIP = EventType("coinflip")

func (ths Coinflip) Type() EventType {
	return COINFLIP
}

var CARDDRAW = EventType("carddraw")

func (ths CardDraw) Type() EventType {
	return CARDDRAW
}

var ATTACHENERGY = EventType("arrachenergy")

func (ths AttachEnergy) Type() EventType {
	return ATTACHENERGY
}
