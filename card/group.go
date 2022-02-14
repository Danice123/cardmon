package card

import (
	"strings"
)

type CardGroup []Card

func (ths *CardGroup) Add(c Card) {
	(*ths) = append((*ths), c)
}

func (ths *CardGroup) Remove(index int) Card {
	element := (*ths)[index]
	(*ths)[index] = (*ths)[len(*ths)-1]
	(*ths) = (*ths)[:len(*ths)-1]
	return element
}

func (ths CardGroup) String() string {
	sb := strings.Builder{}
	for _, c := range ths {
		sb.WriteString(c.String())
	}
	return sb.String()
}

func (ths CardGroup) Slice() []string {
	s := []string{}
	for _, c := range ths {
		s = append(s, c.String())
	}
	return s
}
