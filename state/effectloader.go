package state

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func LoadEffect(id string, param map[string]interface{}) Effect {
	var e Effect
	b, err := yaml.Marshal(param)
	if err != nil {
		panic(err)
	}

	switch id {
	case "coinflip":
		fmt.Printf("%s", string(b))
		var coinflip Coinflip
		err = yaml.Unmarshal(b, &coinflip)
		e = coinflip
	case "protect":
		var protect Protect
		err = yaml.Unmarshal(b, &protect)
		e = protect
	default:
		panic("No effect by that id")
	}
	if err != nil {
		panic(err)
	}
	return e
}
