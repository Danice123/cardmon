package card

type ErrorCard struct{}

func (ErrorCard) Id() string {
	return "error"
}

func (ErrorCard) CardType() CardType {
	return CardType("ERROR")
}

func (ErrorCard) String() string {
	return "error"
}

func (ErrorCard) instanceWithId(string) Card {
	panic("instance of error is invalid")
}
