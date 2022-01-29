package card

import (
	"fmt"
	"strings"
)

type MonsterCard struct {
	Name       string
	Level      int
	Stage      int
	HP         int
	Type       Type
	Weakness   Type
	Resistance Type
	Retreat    int
	Attacks    []MonsterAttack
}

func (ths MonsterCard) CardType() CardType {
	return Monster
}

func (ths MonsterCard) String() string {
	return fmt.Sprintf("%s  LV. %d", ths.Name, ths.Level)
}

func (ths MonsterCard) Text() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%s  %s\tLV. %d\tHP: %d\tTYPE: %s\n", StageToString(ths.Stage), ths.Name, ths.Level, ths.HP, ths.Type))

	for _, attack := range ths.Attacks {
		sb.WriteString(attack.String())
	}
	sb.WriteString(fmt.Sprintf("Weakness: %s\tResistance: %s\tRetreat Cost: %d\n", ths.Weakness, ths.Resistance, ths.Retreat))

	return sb.String()
}

func StageToString(stage int) string {
	if stage == 1 {
		return "Basic"
	} else {
		return fmt.Sprintf("Stage %d", stage-1)
	}
}

type MonsterAttack struct {
	Name   string
	Cost   map[Type]int
	Damage int
}

func (ths MonsterAttack) CheckCost(energyCards CardGroup) bool {
	emap := map[Type]int{}
	for _, c := range energyCards {
		emap[c.(EnergyCard).Type]++
	}

	total := 0
	for _, c := range ths.Cost {
		total += c
	}

	if len(energyCards) < total {
		return false
	}

	ok := true
	for t, c := range ths.Cost {
		if t == COLORLESS {
			continue
		}
		if emap[t] < c {
			ok = false
			break
		}
	}
	return ok
}

func (ths MonsterAttack) String() string {
	return fmt.Sprintf("%s\t%d\tCost: %s\n", ths.Name, ths.Damage, CostToString(ths.Cost))
}

func CostToString(cost map[Type]int) string {
	var sb strings.Builder
	for t, c := range cost {
		sb.WriteString(fmt.Sprintf("%s:%d ", t, c))
	}
	return sb.String()
}
