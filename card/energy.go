package card

import "fmt"

type EnergyCard struct {
	Name string
	Type Type
}

func (ths EnergyCard) ShortText() string {
	return ths.Name
}

func (ths EnergyCard) Text() {
	fmt.Printf("%s\n", ths.Name)
}
