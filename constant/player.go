package constant

type Player string

const Player1 = Player("one")
const Player2 = Player("two")

func OtherPlayer(p Player) Player {
	if p == Player1 {
		return Player2
	} else {
		return Player1
	}
}
