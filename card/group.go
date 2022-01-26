package card

import "fmt"

type CardGroup []Card

func (ths CardGroup) Show() {
	for _, c := range ths {
		fmt.Println(c.ShortText())
	}
}

func (ths *CardGroup) Remove(index int) Card {
	element := (*ths)[index]
	(*ths)[index] = (*ths)[len(*ths)-1]
	(*ths) = (*ths)[:len(*ths)-1]
	return element
}

func (ths CardGroup) ShowIndexed() {
	for i, c := range ths {
		fmt.Printf("%d: %s\n", i, c.ShortText())
	}
}
