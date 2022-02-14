package main

import (
	"fmt"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/machine"
	"github.com/manifoldco/promptui"
)

type HumanConsole struct {
	player  constant.Player
	color   string
	machine *machine.GameMachine
}

func (ths *HumanConsole) Alert(message string) {
	fmt.Printf("%s%s\033[0m\n", ths.color, message)
}

func (ths *HumanConsole) AskCardFromHand(message string, cancelable bool) (int, bool) {
	hand := ths.machine.Current.Players[ths.player].Hand.Slice()
	if cancelable {
		hand = append(hand, "Cancel")
	}

	fmt.Printf("%s%s\033[0m\n", ths.color, message)
	prompt := promptui.Select{
		Label: "Choose a card from your hand",
		Items: hand,
		Size:  10,
	}

	choice, _, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	if choice == len(ths.machine.Current.Players[ths.player].Hand) {
		return -1, false
	}
	return choice, true
}

func (ths *HumanConsole) AskTargetMonster(message string, showActive bool, cancelable bool) (int, bool) {
	choices := []string{}
	if showActive {
		choices = append(choices, fmt.Sprintf("(Active) %s", ths.machine.Current.Players[ths.player].Active.Card.String()))
	}
	for _, c := range ths.machine.Current.Players[ths.player].Bench {
		choices = append(choices, c.Card.String())
	}
	if cancelable {
		choices = append(choices, "Cancel")
	}

	fmt.Printf("%s%s\033[0m\n", ths.color, message)
	prompt := promptui.Select{
		Label: "Choose a target monster",
		Items: choices,
		Size:  6,
	}

	choice, _, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	if cancelable && choice == len(choices)-1 {
		return -1, false
	}
	return choice, true
}

func (ths *HumanConsole) AskForAction() (string, int) {
	fmt.Printf("%s%s's Turn: ", ths.color, ths.player)
	for {
		var err error
		var choice int
		prompt := promptui.Select{
			Label: "Actions",
			Items: []string{"Hand", "Attack", "Retreat", "Pass"},
		}
		if choice, _, err = prompt.Run(); err != nil {
			panic(err)
		}

		switch choice {
		case 0:
			c, ok := ths.AskCardFromHand("Which card", true)
			if ok {
				return "Hand", c
			}
		case 1:
			attacks := []string{}
			for _, attack := range ths.machine.Current.Players[ths.player].Active.Card.(card.MonsterCard).Attacks {
				attacks = append(attacks, attack.String())
			}
			attacks = append(attacks, "Cancel")

			prompt := promptui.Select{
				Label: "Attacks",
				Items: attacks,
			}
			if choice, _, err = prompt.Run(); err != nil {
				panic(err)
			}

			if choice != len(attacks)-1 {
				return "Attack", choice
			}
		case 2:
			c, ok := ths.AskTargetMonster("Which monster to switch to?", false, true)
			if ok {
				return "Retreat", c
			}
		case 3:
			return "Pass", 0
		}
	}
}
