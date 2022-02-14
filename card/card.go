package card

type Card interface {
	Id() string
	CardType() CardType
	String() string

	instanceWithId(string) Card
}
