package card

import (
	"fmt"
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

func (ths CardGroup) Show() {
	for _, c := range ths {
		fmt.Println(c.String())
	}
}

func (ths CardGroup) ShowIndexed() {
	for i, c := range ths {
		fmt.Printf("%d: %s\n", i+1, c.String())
	}
}
