package discord

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jmsheff/discord-checkers/logic"
)

// Adds the reactions to give the user the ability to select a piece
func addSelectReactions(s *discordgo.Session, c string, m string, g *logic.Game) {
	// Adds all the Y coordinate selections
	for i, e := range ySlice {
		// Gives only the rows that have a piece in them
		if strings.ContainsAny(g.Board[i*4:i*4+4], strconv.FormatUint(uint64(g.Turn), 10)+strconv.FormatUint(uint64(g.Turn+2), 10)) {
			s.MessageReactionAdd(c, m, e)
		}
	}

	// Adds all the x coordinate selections
	for _, e := range xSlice {
		// Because of rate limits around reactions, it's not suitable to dynamically add x reactions like we do with y. Therefore we add them all
		s.MessageReactionAdd(c, m, e)
	}

	s.MessageReactionAdd(c, m, "✅")
}

// Gets an Y coordinate from a selection emoji
func getYFromReaction(e *discordgo.Emoji) uint8 {
	for i, y := range ySlice {
		if e.Name == y {
			return uint8(i + 1)
		}
	}
	return 0
}

// Gets an X coordinate from a selection emoji
func getXFromReaction(e *discordgo.Emoji) uint8 {
	for i, x := range xSlice {
		if e.Name == x {
			return uint8(i + 1)
		}
	}
	return 0
}

// Selects a piece and shows the moves on the board
func selectPiece(s *discordgo.Session, c string, m string, opponentID string, game *logic.Game, square *logic.Square, moves *[]logic.Move, jumpsOnly bool) {
	// Makes a board with moves and a slice of reactions to put on the message
	board := []rune(game.Board)
	var reactions []string
	for i, move := range *moves {
		if move.Possible {
			board[move.S.Index] = directionSlice[i]
			reactions = append(reactions, movesSlice[i])
		}
	}

	// Select the piece
	game.Selected = square.Index

	// Send the board with the moves on it
	gamemsg, err := s.ChannelMessageSendEmbed(c, gameEmbed(s, "move", opponentID, game, string(board), false))
	if err != nil {
		return
	}

	// Add the move reactions
	if !jumpsOnly { // Don't allow cancelling on double jumps
		s.MessageReactionAdd(gamemsg.ChannelID, gamemsg.ID, "❌")
	}
	for _, e := range reactions {
		s.MessageReactionAdd(gamemsg.ChannelID, gamemsg.ID, e)
	}

	s.ChannelMessageDelete(c, m) // The reason we delete instead of edit is to get around not being able to clear reactions
}

// Handles all selection related reactions
func selectReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd, m *discordgo.Message, user *discordgo.User, gameString string) {

	if r.Emoji.Name == "✅" { // Only verfiy if the user is confirming it
		// Prevent reaction spam
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, s.State.User.ID)
		defer s.MessageReactionAdd(r.ChannelID, r.MessageID, r.Emoji.Name) // Add it back after if there is an error

		// Get the users reactions that the bot has also put on
		var userReactions []*discordgo.MessageReactions
		botReactions := getBotReactions(m.Reactions)

		for _, botReaction := range botReactions {
			if botReaction.Count > 1 {
				if botReaction.Emoji.Name != "✅" {
					userReactions = append(userReactions, botReaction)
				}
			}
		}

		// Make sure right length
		if len(userReactions) != 2 {
			// Send error message to warn them of their mistake
			s.ChannelMessageSend(r.ChannelID, errorMessage("Invalid reactions", "Ensure you have reacted with 1 letter and 1 number. Adjust your reactions and try again."))
			return
		}

		x1 := getXFromReaction(userReactions[0].Emoji)
		y1 := getYFromReaction(userReactions[0].Emoji)
		x2 := getXFromReaction(userReactions[1].Emoji)
		y2 := getYFromReaction(userReactions[1].Emoji)

		var x uint8
		var y uint8
		// If the first reaction is the x and the second reaction is the y
		if x1 != 0 && y2 != 0 {
			x = x1
			y = y2
		} else if y1 != 0 && x2 != 0 { // If the first reaction is the y and the second reaction is the x
			y = y1
			x = x2
		} else {
			// Most likely reacted with 2 numbers or 2 letters
			s.ChannelMessageSend(r.ChannelID, errorMessage("Invalid reactions", "Ensure you reacted with 1 letter and 1 number. Adjust your reactions and try again."))
			return
		}

		// Get game
		opponentID, game, err := ParseGame(gameString)
		if err != nil {
			return
		}

		// Get selection
		square, _ := logic.SquareAtCoords(x, y, &game)
		moves, err := square.GetAvailableMoves(&game, false)
		if err != nil {
			s.ChannelMessageSend(r.ChannelID, errorMessage(err.Error(), "Adjust your reactions and try again."))
			return
		}

		// If all is good, then we can get the available moves
		selectPiece(s, r.ChannelID, r.MessageID, opponentID, &game, &square, &moves, false)
	}
}
