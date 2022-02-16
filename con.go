package main

import (
	"fmt"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/constant"
	"github.com/Danice123/cardmon/machine"
	"github.com/Danice123/cardmon/state"
	"github.com/manifoldco/promptui"
)

type HumanConsole struct {
	player  constant.Player
	color   string
	machine *machine.GameMachine
}

func (ths *HumanConsole) Handle(event machine.Event) {
	switch event.Type() {
	case machine.COINFLIP:
		fmt.Printf("%s%s\nThe result was %s\033[0m\n", ths.color, event.(machine.Coinflip).Message, event.(machine.Coinflip).Outcome)
	case machine.CARDDRAW:
		fmt.Printf("%sYou drew the card: %s\033[0m\n", ths.color, event.(machine.CardDraw).Card.String())
	case machine.ATTACHENERGY:
		if event.(machine.AttachEnergy).Player == ths.player {
			fmt.Printf("%sYou attached the %s to %s\033[0m\n",
				ths.color,
				event.(machine.AttachEnergy).Energy.String(),
				event.(machine.AttachEnergy).Target.String())
		} else {
			fmt.Printf("%sOpponent attached the %s to %s\033[0m\n",
				ths.color,
				event.(machine.AttachEnergy).Energy.String(),
				event.(machine.AttachEnergy).Target.String())
		}
	}
}

func (ths *HumanConsole) Alert(message string) {
	fmt.Printf("%s%s\033[0m\n", ths.color, message)
}

func (ths *HumanConsole) AskCardFromHand(message string, cancelable bool) (string, bool) {
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
		return "", false
	}
	return ths.machine.Current.Players[ths.player].Hand[choice].Id(), true
}

func (ths *HumanConsole) AskTargetMonster(message string, showActive bool, cancelable bool) (string, bool) {
	choices := []string{}
	ids := []string{}
	if showActive {
		choices = append(choices, fmt.Sprintf("(Active) %s", ths.machine.Current.Players[ths.player].Active.Card.String()))
		ids = append(ids, ths.machine.Current.Players[ths.player].Active.Card.Id())
	}
	for _, c := range ths.machine.Current.Players[ths.player].Bench {
		choices = append(choices, c.Card.String())
		ids = append(ids, c.Card.Id())
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
		return "", false
	}
	return ids[choice], true
}

func (ths *HumanConsole) AskForAction() (string, string) {
	fmt.Printf("%s%s's Turn:\033[0m\n", ths.color, ths.player)
	PrintGame(ths.machine.Current)
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
			ids := []string{}
			for _, attack := range ths.machine.Current.Players[ths.player].Active.Card.(card.MonsterCard).Attacks {
				attacks = append(attacks, attack.String())
				ids = append(ids, attack.Name)
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
				return "Attack", ids[choice]
			}
		case 2:
			c, ok := ths.AskTargetMonster("Which monster to switch to?", false, true)
			if ok {
				return "Retreat", c
			}
		case 3:
			return "Pass", ""
		}
	}
}

func PrintGame(gs state.Gamestate) {
	fmt.Printf("#\t")
	for _, c := range gs.Players[constant.Player2].Bench {
		fmt.Printf("%s\t", c.Card.String())
	}
	fmt.Print("\n")
	fmt.Printf("#\n#\t%s\n#\n", gs.Players[constant.Player2].Active.String())
	fmt.Println("######################################################")
	fmt.Printf("#\n#\t%s\n#\n", gs.Players[constant.Player1].Active.String())
	fmt.Printf("#\t")
	for _, c := range gs.Players[constant.Player1].Bench {
		fmt.Printf("%s\t", c.Card.String())
	}
	fmt.Print("\n")
}
