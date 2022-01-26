package card

import "math/rand"

type CardStack CardGroup

func (ths *CardStack) Push(card Card) {
	*ths = append(*ths, card)
}

func (ths *CardStack) Pop() (Card, bool) {
	if len(*ths) == 0 {
		return nil, false
	} else {
		index := len(*ths) - 1
		element := (*ths)[index]
		*ths = (*ths)[:index]
		return element, true
	}
}

func (ths *CardStack) PopX(n int) []Card {
	cards := []Card{}
	index := len(*ths) - n
	if index < 0 {
		index = 0
	}
	for i := index; i < len(*ths); i++ {
		cards = append(cards, (*ths)[i])
	}
	*ths = (*ths)[:index]
	return cards
}

func (ths CardStack) Shuffle() {
	rand.Shuffle(len(ths), func(i, j int) {
		ths[i], ths[j] = ths[j], ths[i]
	})
}
