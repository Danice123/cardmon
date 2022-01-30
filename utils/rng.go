package utils

import "math/rand"

func Coinflip() bool {
	if rand.Intn(100) > 50 {
		return true
	} else {
		return false
	}
}
