package logic

import (
	"errors"
	"math"
	"strconv"
)

// Basically acts as an enchanced index for the board
type Square struct {
	X     uint8 // X coordinate of the square
	Y     uint8 // Y coordinate of the square
	Index uint8 // Index of the square on the board
	Piece uint8 // The piece located at the square
}

// A direction used for getting moves and squares at a direction relative to the piece
type Direction struct {
	Y      string // North or south
	X      string // East or west
	Offset int8   // Offset for the index on an even row
}

// Gives the directions as an enum
var NORTHWEST Direction = Direction{Y: "north", X: "west", Offset: -4}
var NORTHEAST Direction = Direction{Y: "north", X: "east", Offset: -3}
var SOUTHWEST Direction = Direction{Y: "south", X: "west", Offset: 4}
var SOUTHEAST Direction = Direction{Y: "south", X: "east", Offset: 5}

// Iterable slice to loop through all directions
var Directions []Direction = []Direction{NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST}

// Gets a square from coordinates
func SquareAtCoords(x uint8, y uint8, game *Game) (Square, error) {
	index := ((y - 1) * 4) + (x - 1)
	piece, err := strconv.ParseUint(string((game.Board)[index]), 10, 0)

	if err != nil {
		return Square{}, errors.New("Error parsing piece")
	}

	return Square{
		X:     x - 1,
		Y:     y - 1,
		Index: index,
		Piece: uint8(piece),
	}, nil
}

// Gets a square at an index
func SquareAtIndex(index uint8, game *Game) (Square, error) {
	x := index % 4
	y := uint8(math.Ceil(float64(index / 4)))
	piece, err := strconv.ParseUint(string(game.Board[index]), 10, 0)

	if err != nil {
		return Square{}, errors.New("Error parsing piece")
	}

	return Square{
		X:     x,
		Y:     y,
		Index: index,
		Piece: uint8(piece),
	}, nil

}

// Checks if a square is empty
func (s Square) IsEmpty() bool {
	return s.Piece == 0
}

// Gets the player returns 1 or 2
func (s Square) Player() uint8 {
	if s.Piece > 2 { // If the piece is a king
		return s.Piece - 2
	} else {
		return s.Piece
	}
}

// Checks if the piece is a king
func (s Square) IsKing() bool {
	return s.Piece > 2
}

// Gets a square at a given direction
func (s Square) SquareAtDirection(direction Direction, game *Game) (Square, error) {
	// Make sure isn't out of range
	if (s.Y == 0 && direction.Y == "north") || (s.Y == 7 && direction.Y == "south") {
		return Square{}, errors.New("Out of bounds")
	}
	// If square is even
	if s.Y%2 == 0 {
		if s.X != 3 || direction.X == "west" {
			return SquareAtIndex(uint8(int8(s.Index)+direction.Offset), game)
		} else {
			return Square{}, errors.New("Out of bounds")
		}
	} else {
		if s.X != 0 || direction.X == "east" {
			return SquareAtIndex(uint8(int8(s.Index)+direction.Offset-1), game)
		} else {
			return Square{}, errors.New("Out of bounds")
		}
	}
}

// Gets a move at a given direction
func (s Square) MoveAtDirection(direction Direction, game *Game, jumpOnly bool) Move {
	squareAtDir, err := s.SquareAtDirection(direction, game)

	// Makes sure its a valid square and that it can only move south if it's a king
	if err == nil && (direction.Y == "north" || s.IsKing()) {
		// Checks for a normal move
		if squareAtDir.IsEmpty() && !jumpOnly {
			return Move{Possible: true, S: squareAtDir}
		} else if !squareAtDir.IsEmpty() && squareAtDir.Player() != game.Turn { // Checks if it is the oppisite player
			if jump, err := squareAtDir.SquareAtDirection(direction, game); err == nil && jump.IsEmpty() { // Check if jump is possible
				return Move{Possible: true, S: jump, Jumped: squareAtDir}
			}
		}
	}
	// Returning a blank move will make the "Possible" field false, making it so we can tell which moves are possible when we iterate
	return Move{}
}

// Gets all moves avaiable in every direction
func (s Square) GetAvailableMoves(game *Game, jumpOnly bool) ([]Move, error) {
	if s.X < 0 || s.X > 3 || s.Y < 0 || s.Y > 7 {
		return []Move{}, errors.New("Invalid coordinates")
	}

	// Validate piece
	if s.IsEmpty() {
		return []Move{}, errors.New("Cannot select blank space")
	}

	if s.Player() != game.Turn {
		return []Move{}, errors.New("Cannot select other players piece")
	}

	// Append each directions move to the list
	var avaliableMoves []Move
	possibleMoves := false
	for _, dir := range Directions {
		move := s.MoveAtDirection(dir, game, jumpOnly)
		if move.Possible {
			possibleMoves = true
		}
		avaliableMoves = append(avaliableMoves, move)
	}

	if !possibleMoves {
		return []Move{}, errors.New("Piece has nowhere to move")
	}

	return avaliableMoves, nil
}
