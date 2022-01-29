package card

type EnergyCard struct {
	Name string
	Type Type
}

func (ths EnergyCard) CardType() CardType {
	return Energy
}

func (ths EnergyCard) String() string {
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

func (ths EnergyCard) Text() string {
	return ths.Name
}
