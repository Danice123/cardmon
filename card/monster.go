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

func (ths MonsterCard) ShortText() string {
	return fmt.Sprintf("%s  LV. %d", ths.Name, ths.Level)
}

func (ths MonsterCard) Text() {
	fmt.Printf("%s  %s\tLV. %d\tHP: %d\tTYPE: %s\n", StageToString(ths.Stage), ths.Name, ths.Level, ths.HP, ths.Type)
	for _, attack := range ths.Attacks {
		fmt.Printf("%s\t%d\tCost: %s\n", attack.Name, attack.Damage, CostToString(attack.Cost))
	}
	fmt.Printf("Weakness: %s\tResistance: %s\tRetreat Cost: %d\n", ths.Weakness, ths.Resistance, ths.Retreat)
}

func StageToString(stage int) string {
	if stage == 1 {
		return "Basic"
	} else {
		return fmt.Sprintf("Stage %d", stage-1)
	}
}

func CostToString(cost map[Type]int) string {
	var sb strings.Builder
	for t, c := range cost {
		sb.WriteString(fmt.Sprintf("%s:%d ", t, c))
	}
	return sb.String()
}

type MonsterAttack struct {
	Name   string
	Cost   map[Type]int
	Damage int
}
