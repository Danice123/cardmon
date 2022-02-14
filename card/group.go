package card

type CardGroup []Card

func (ths CardGroup) indexOf(id string) (int, bool) {
	index := -1
	for i, c := range ths {
		if c.Id() == id {
			index = i
			break
		}
	}
	return index, index != -1
}

func (ths CardGroup) Get(id string) (Card, bool) {
	if index, ok := ths.indexOf(id); ok {
		return ths[index], true
	}
	return ErrorCard{}, false
}

func (ths *CardGroup) Add(c Card) {
	(*ths) = append((*ths), c)
}

func (ths *CardGroup) Remove(id string) (Card, bool) {
	if index, ok := ths.indexOf(id); ok {
		element := (*ths)[index]
		(*ths)[index] = (*ths)[len(*ths)-1]
		(*ths) = (*ths)[:len(*ths)-1]
		return element, true
	}
	return ErrorCard{}, false
}

func (ths CardGroup) Slice() []string {
	s := []string{}
	for _, c := range ths {
		s = append(s, c.String())
	}
	return s
}
