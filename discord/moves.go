package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/jmsheff/discord-checkers/logic"
)

// Gets the given move from the reaction
func getMoveFromReaction(e *discordgo.Emoji, square *logic.Square, game *logic.Game) (logic.Move, error) {
	var moveTo logic.Move
	for i, j := range movesSlice {
		if e.Name == j {
			if move := square.MoveAtDirection(logic.Directions[i], game, false); move.Possible {
				moveTo = move
				break
			} else {
				return moveTo, errors.New("Move not possible")
			}
		}
	}

	return moveTo, nil
}

// Handles all move related reactions
func moveReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd, m *discordgo.Message, user *discordgo.User, gameString string) {
	opponentID, game, err := ParseGame(gameString)
	// Allows there to only be one reaction present at a time to prevent reaction spam
	if hasOtherReactionsBesides(r.Emoji.Name, m.Reactions) {
		return
	}

	// Deselect
	if r.Emoji.Name == "❌" {
		game.Selected = 0
		if err != nil {
			s.ChannelMessageEdit(r.ChannelID, r.MessageID, errorMessage("Bot error", "Could not deselect piece"))
		}

		gamemsg, _ := s.ChannelMessageSendEmbed(r.ChannelID, gameEmbed(s, "select", opponentID, &game, game.Board, false))
		s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		addSelectReactions(s, r.ChannelID, gamemsg.ID, &game)

		return
	}
	if err != nil {
		s.ChannelMessageEdit(r.ChannelID, r.MessageID, errorMessage("Error getting game", "This game has most likely expired or doesn't exist."))
		return
	}

	// These both will not have errors unless there is flawed logic in selection
	square, _ := logic.SquareAtIndex(game.Selected, &game)
	move, _ := getMoveFromReaction(&r.Emoji, &square, &game)

	logic.MovePiece(square, move, &game.Board)
	updatedSquare, _ := logic.SquareAtIndex(move.S.Index, &game)

	if (move.Jumped != logic.Square{}) { // If move was a jump
		// Check for win
		if p1score, p2score := logic.GetScore(&game); p1score == 12 || p2score == 12 {
			// err := cache.DeleteGame(identifier)
			if err != nil {
				s.ChannelMessageSend(r.ChannelID, errorMessage("Bot error", "Could not delete game."))
				return
			}

			loser, _ := s.User(opponentID)
			loserDM, _ := s.UserChannelCreate(opponentID)
			s.ChannelMessageSendEmbed(r.ChannelID, &discordgo.MessageEmbed{
				Title:       "🎉 YOU WIN!!! 🏆",
				Description: "Congratulations! You won the game against " + formatUser(loser),
				Color:       c_GREEN,
			})

			s.ChannelMessageSendEmbed(loserDM.ID, &discordgo.MessageEmbed{
				Title:       "❌ You lost. ❌",
				Description: "You lost the game against " + formatUser(user) + ". Better luck next time!",
				Color:       c_RED,
			})
			return
		}

		// Check for double jump
		if doubleJumps, err := updatedSquare.GetAvailableMoves(&game, true); err == nil {
			// Selects the piece at the updated location and provides only double jump moves
			selectPiece(s, r.ChannelID, opponentID, &game, &move.S, &doubleJumps, true)
			return
		}
	}

	// Confirm with the current player that their move went through
	s.ChannelMessageSend(r.ChannelID, successMessage("Move sent!", "Wait here for them to make their move."))
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, gameEmbed(s, "", opponentID, &game, game.Board, true)) // Keep a record of the move

	// If no double jump swap turn
	err = logic.SwapTurn(&game)
	if err != nil {
		s.ChannelMessageEdit(r.ChannelID, r.MessageID, errorMessage("Bot error", "Could not swap turn"))
		return
	}

	// Send game to opponent for their move
	opponentDM, err := s.UserChannelCreate(opponentID)
	opponentMsg, err := s.ChannelMessageSendEmbed(opponentDM.ID, gameEmbed(s, "select", r.UserID, &game, game.Board, false))
	addSelectReactions(s, opponentMsg.ChannelID, opponentMsg.ID, &game)
}
