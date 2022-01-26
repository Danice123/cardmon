package card

import (
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Library map[string]map[int]Card

type CardType string

const Monster = CardType("Monster")
const Trainer = CardType("Trainer")
const Energy = CardType("Energy")

type GenericCard struct {
	Card   CardType
	Series string
	Number int
}

func LoadLibrary() Library {
	library := Library{}

	err := filepath.WalkDir("assets/cards", func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			d, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var gc GenericCard
			err = yaml.Unmarshal(d, &gc)
			if err != nil {
				return err
			}

			if _, ok := library[gc.Series]; !ok {
				library[gc.Series] = map[int]Card{}
			}

			switch gc.Card {
			case Monster:
				var c MonsterCard
				err = yaml.Unmarshal(d, &c)
				library[gc.Series][gc.Number] = c
			case Trainer:
				// TODO
			case Energy:
				var c EnergyCard
				err = yaml.Unmarshal(d, &c)
				library[gc.Series][gc.Number] = c
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return library
}
