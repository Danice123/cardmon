package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Danice123/cardmon/card"
	"github.com/Danice123/cardmon/deck"
	"github.com/Danice123/cardmon/machine"
	"github.com/Danice123/cardmon/state"
)

var console = bufio.NewScanner(os.Stdin)

func consoleReadInt(bound int) int {
	for {
		console.Scan()
		i, err := strconv.Atoi(console.Text())
		if err == nil && i > 0 && i <= bound {
			return i
		}
		fmt.Println("Try again: ")
	}
}

type HumanConsole struct {
	player  state.Player
	color   string
	machine *machine.GameMachine
}

func (ths *HumanConsole) Alert(message string) {
	fmt.Printf("%s%s: %s\033[0m\n", ths.color, ths.player, message)
}

func (ths *HumanConsole) AskCardFromHand(message string, cancelable bool) (int, bool) {
	fmt.Printf("%s%s: %s\033[0m\n", ths.color, ths.player, message)
	ths.machine.Current.Players[ths.player].Hand.ShowIndexed()
	if cancelable {
		fmt.Printf("%d: Cancel\n", len(ths.machine.Current.Players[ths.player].Hand)+1)
	}
	fmt.Print("Choice? ")

	if cancelable {
		choice := consoleReadInt(len(ths.machine.Current.Players[ths.player].Hand) + 1)
		if choice == len(ths.machine.Current.Players[ths.player].Hand)+1 {
			return -1, false
		}
		return choice - 1, true
	} else {
		return consoleReadInt(len(ths.machine.Current.Players[ths.player].Hand)) - 1, true
	}
}

func (ths *HumanConsole) AskTargetMonster(message string, cancelable bool) (int, bool) {
	fmt.Printf("%s%s: %s\033[0m\n", ths.color, ths.player, message)
	fmt.Printf("1 (Active): %s\n", ths.machine.Current.Players[ths.player].Active.Card.String())
	for i, c := range ths.machine.Current.Players[ths.player].Bench {
		fmt.Printf("%d (Active): %s\n", i+2, c.Card.String())
	}
	if cancelable {
		fmt.Printf("%d: Cancel\n", len(ths.machine.Current.Players[ths.player].Bench)+2)
	}
	fmt.Print("Choice? ")

	if cancelable {
		choice := consoleReadInt(len(ths.machine.Current.Players[ths.player].Bench) + 2)
		if choice == len(ths.machine.Current.Players[ths.player].Bench)+2 {
			return -1, false
		}
		return choice - 1, true
	} else {
		return consoleReadInt(len(ths.machine.Current.Players[ths.player].Bench)+1) - 1, true
	}
}

func (ths *HumanConsole) AskTargetBenchMonster(message string, cancelable bool) (int, bool) {
	fmt.Printf("%s%s: %s\033[0m\n", ths.color, ths.player, message)
	for i, c := range ths.machine.Current.Players[ths.player].Bench {
		fmt.Printf("%d: %s\n", i+1, c.Card.String())
	}
	if cancelable {
		fmt.Printf("%d: Cancel\n", len(ths.machine.Current.Players[ths.player].Bench)+1)
	}
	fmt.Print("Choice? ")

	if cancelable {
		choice := consoleReadInt(len(ths.machine.Current.Players[ths.player].Bench) + 1)
		if choice == len(ths.machine.Current.Players[ths.player].Bench)+1 {
			return -1, false
		}
		return choice - 1, true
	} else {
		return consoleReadInt(len(ths.machine.Current.Players[ths.player].Bench)) - 1, true
	}
}

func (ths *HumanConsole) AskForAction() (string, int) {
	fmt.Printf("%s%s: ", ths.color, ths.player)
	handSize := len(ths.machine.Current.Players[ths.player].Hand)
	for {
		fmt.Print("Actions:\n1. Hand\t2. Attack\n3. Retreat\t4. Pass\n\033[0m")
		fmt.Print("Choice? ")
		switch consoleReadInt(4) {
		case 1:
			ths.machine.Current.Players[ths.player].Hand.ShowIndexed()
			fmt.Printf("%d. Back\n", handSize+1)
			c := consoleReadInt(handSize+1) - 1
			if c != handSize {
				return "Hand", c
			}
		case 2:
			attackSize := len(ths.machine.Current.Players[ths.player].Active.Card.(card.MonsterCard).Attacks)
			for i, attack := range ths.machine.Current.Players[ths.player].Active.Card.(card.MonsterCard).Attacks {
				fmt.Printf("%d: %s", i+1, attack.String())
			}
			fmt.Printf("%d. Back\n", attackSize+1)
			fmt.Print("Choice? ")
			c := consoleReadInt(attackSize+1) - 1
			if c != attackSize {
				return "Attack", c
			}
		case 3:
			c, ok := ths.AskTargetBenchMonster("Which monster to switch to?", true)
			if ok {
				return "Retreat", c
			}
		case 4:
			return "Pass", 0
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	lib := card.LoadLibrary()
	gs := state.NewGame(deck.BuildDeck("assets/decks/test.yaml", lib), deck.BuildDeck("assets/decks/test.yaml", lib))

	g := &machine.GameMachine{
		Current: gs,
	}
	g.RegisterHandler(state.Player1, &HumanConsole{
		player:  state.Player1,
		color:   "\033[31m",
		machine: g,
	})
	g.RegisterHandler(state.Player2, &HumanConsole{
		player:  state.Player2,
		color:   "\033[34m",
		machine: g,
	})
	g.Start(gs)
}
