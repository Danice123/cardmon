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

type HumanConsole struct {
	player  state.Player
	color   string
	machine *machine.GameMachine
}

func (ths *HumanConsole) Alert(message string) {
	fmt.Printf("%s%s: %s\033[0m\n", ths.color, ths.player, message)
}

func (ths *HumanConsole) AskCardFromHand(message string) int {
	fmt.Printf("%s%s: %s\n", ths.color, ths.player, message)
	ths.machine.Current.Players[ths.player].Hand.ShowIndexed()
	fmt.Print("Choice? \033[0m")
	console.Scan()
	c, err := strconv.Atoi(console.Text())
	if err != nil {
		panic(err)
	}
	return c
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
