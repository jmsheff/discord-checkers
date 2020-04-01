package logic

type Game struct {
	Selected uint8  // The index of the selected piece
	Turn     uint8  // Which players turn it is(1 or 2)
	Board    string // The board represented as a string
}
