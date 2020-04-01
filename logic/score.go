package logic

import (
	"strings"
)

// Gets the score of a game
func GetScore(g *Game) (int, int) {
	// Score is in pieces captured so that's why it's a bit counter intuitive
	return 12 - (strings.Count(g.Board, "2") + strings.Count(g.Board, "4")), 12 - (strings.Count(g.Board, "1") + strings.Count(g.Board, "3"))
}
