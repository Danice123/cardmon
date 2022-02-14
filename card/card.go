package card

type Card interface {
	CardType() CardType
	String() string
}
