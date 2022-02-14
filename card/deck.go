package card

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type DeckFile struct {
	Cards map[string]int
}

func BuildDeck(pid string, path string, library Library) CardStack {
	d, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var df DeckFile
	err = yaml.Unmarshal(d, &df)
	if err != nil {
		panic(err)
	}

	deck := CardStack{}
	for card, n := range df.Cards {
		s := strings.Split(card, "_")
		num, err := strconv.Atoi(s[1])
		if err != nil {
			panic(err)
		}
		c := library[s[0]][num]
		for i := 0; i < n; i++ {
			deck = append(deck, c.instanceWithId(fmt.Sprintf("%s_%d", pid, i)))
		}
	}

	return deck
}
