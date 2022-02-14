package card

import "fmt"

type EnergyCard struct {
	id   string
	Name string
	Type Type
}

func (ths EnergyCard) instanceWithId(did string) Card {
	ths.id = fmt.Sprintf("%s_%s", ths.Name, did)
	return ths
}

func (ths EnergyCard) Id() string {
	if ths.id == "" {
		panic("Uninitialized card!")
	}
	return ths.id
}

func (ths EnergyCard) CardType() CardType {
	return Energy
}

func (ths EnergyCard) Pprint() string {
	switch ths.Type {
	case GRASS:
		return "\033[32mE\033[0m"
	case FIRE:
		return "\033[31mE\033[0m"
	case WATER:
		return "\033[34mE\033[0m"
	case ELECTRIC:
		return "\033[33mE\033[0m"
	case FIGHTING:
		return "\033[7;31mE\033[0m"
	case PSYCHIC:
		return "\033[35mE\033[0m"
	default:
		return "E"
	}
}

func (ths EnergyCard) String() string {
	return ths.Name
}
