package state

import (
	"fmt"

	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/utils"
)

type StatusEffect interface {
	OnTurnChange(gs Gamestate, p constant.Player) (Gamestate, []Event)
	CanMonsterAttack() (string, bool)
	CanMonsterRetreat() (string, bool)
}

type StatusList []StatusEffect

func (ths StatusList) indexOf(status StatusEffect) (int, bool) {
	index := -1
	for i, s := range ths {
		if s == status {
			index = i
			break
		}
	}
	return index, index != -1
}

func (ths *StatusList) Add(status StatusEffect) {
	(*ths) = append((*ths), status)
}

func (ths *StatusList) Remove(status StatusEffect) (StatusEffect, bool) {
	if index, ok := ths.indexOf(status); ok {
		element := (*ths)[index]
		(*ths)[index] = (*ths)[len(*ths)-1]
		(*ths) = (*ths)[:len(*ths)-1]
		return element, true
	}
	return nil, false
}

type SleepStatus struct{}

func (ths SleepStatus) OnTurnChange(gs Gamestate, p constant.Player) (Gamestate, []Event) {
	outcome := utils.Coinflip()
	if outcome {
		state := gs.Players[p]
		state.Active.Statuses.Remove(ths)
		gs.Players[p] = state
	}
	return gs, []Event{ECoinflip{Message: fmt.Sprintf("If head, %s wakes up", gs.Players[p].Active.Card.String()), Outcome: outcome}}
}

func (ths SleepStatus) CanMonsterAttack() (string, bool) {
	return "Monster is asleep!", false
}

func (ths SleepStatus) CanMonsterRetreat() (string, bool) {
	return "Monster is asleep!", false
}
