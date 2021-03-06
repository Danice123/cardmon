package card

import (
	"fmt"
	"strings"
)

type MonsterCard struct {
	id         string
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

func (ths MonsterCard) instanceWithId(did string) Card {
	ths.id = fmt.Sprintf("%s_%s", ths.Name, did)
	return ths
}

func (ths MonsterCard) Id() string {
	if ths.id == "" {
		panic("Uninitialized card!")
	}
	return ths.id
}

func (ths MonsterCard) CardType() CardType {
	return Monster
}

func (ths MonsterCard) String() string {
	return fmt.Sprintf("%s  LV. %d", ths.Name, ths.Level)
}

func (ths MonsterCard) GetAttack(aid string) (MonsterAttack, bool) {
	index := -1
	for i, attack := range ths.Attacks {
		if attack.Name == aid {
			index = i
			break
		}
	}
	if index == -1 {
		return MonsterAttack{}, false
	}
	return ths.Attacks[index], true

}

type MonsterAttack struct {
	Name              string
	Cost              map[Type]int
	Description       string
	DamageDescription string
	Effects           []Effect
}

func CostToString(cost map[Type]int) string {
	var sb strings.Builder
	for t, c := range cost {
		sb.WriteString(fmt.Sprintf("%s:%d ", t, c))
	}
	return sb.String()
}

func (ths MonsterAttack) String() string {
	return fmt.Sprintf("%s\t%s\tCost: %s", ths.Name, ths.DamageDescription, CostToString(ths.Cost))
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

type Effect struct {
	Id         string
	Parameters map[string]interface{}
}
