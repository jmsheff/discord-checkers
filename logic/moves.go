package logic

import (
	"errors"
	"strconv"
)

// Represents a move for a piece
type Move struct {
	Possible bool   // If the move is possible
	S        Square // The Square to move the piece to
	Jumped   Square // The square of the jumped piece if the move was a jump
}

// Moves a piece and updates the board
func MovePiece(s Square, m Move, board *string) {
	boardRunes := []rune(*board)
	if m.S.Y == 0 && s.Piece < 3 {
		// King the piece if it is in the right row and isn't already kinged
		boardRunes[m.S.Index] = rune(strconv.FormatUint(uint64(s.Piece+2), 10)[0])
	} else {
		// Copy piece to new position
		boardRunes[m.S.Index] = rune(strconv.FormatUint(uint64(s.Piece), 10)[0])
	}
	// Remove old piece
	boardRunes[s.Index] = rune('0')
	// If the move was a jump remove jumped piece
	if (m.Jumped != Square{}) {
		boardRunes[m.Jumped.Index] = rune('0')
	}

	*board = string(boardRunes)
}

// Swaps the turn for a game
func SwapTurn(game *Game) error {
	// Removes the selection
	game.Selected = 0

	// Toggles the turn
	if game.Turn == 1 {
		game.Turn = 2
	} else if game.Turn == 2 {
		game.Turn = 1
	} else {
		return errors.New("Invalid turn")
	}

	// Flips the board perspective
	boardRunes := []rune(game.Board)
	for i, j := 0, len(boardRunes)-1; i < j; i, j = i+1, j-1 {
		boardRunes[i], boardRunes[j] = boardRunes[j], boardRunes[i]
	}
	game.Board = string(boardRunes)

	return nil
}
