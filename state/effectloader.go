package state

import (
	"gopkg.in/yaml.v3"
)

func LoadEffect(id string, param map[string]interface{}) Effect {
	var e Effect
	b, err := yaml.Marshal(param)
	if err != nil {
		panic(err)
	}

	switch id {
	case "damage":
		var effect Damage
		err = yaml.Unmarshal(b, &effect)
		e = effect
	case "coinflip":
		var effect Coinflip
		err = yaml.Unmarshal(b, &effect)
		e = effect
	case "sleep":
		var effect SleepEffect
		err = yaml.Unmarshal(b, &effect)
		e = effect
	case "protect":
		var effect Protect
		err = yaml.Unmarshal(b, &effect)
		e = effect
	case "selfdamage":
		var effect Selfdamage
		err = yaml.Unmarshal(b, &effect)
		e = effect
	case "metronome":
		var effect Metronome
		err = yaml.Unmarshal(b, &effect)
		e = effect
	default:
		panic("No effect by that id")
	}
	if err != nil {
		panic(err)
	}
	return e
}
