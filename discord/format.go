package discord

import (
	"math"
	"strings"

	"github.com/jmsheff/discord-checkers/logic"

	"github.com/bwmarrin/discordgo"
)

// Some helper "Enums"
var ySlice []string = []string{"🇦", "🇧", "🇨", "🇩", "🇪", "🇫", "🇬", "🇭"}
var xSlice []string = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣"}
var directionSlice []rune = []rune{'N', 'n', 'S', 's'}
var movesSlice []string = []string{"↖️", "↗️", "↙️", "↘️"}

// Make the board into emojis
func formatBoard(board *string) string {
	var formatted []string

	for i, j := range *board {
		// The current row
		row := int(math.Ceil(float64(i / 4)))
		var e string

		// If we are at the start of the row
		if (i+1)%4 == 1 {
			formatted = append(formatted, ySlice[row])
		}

		// Get what the emoji is based on the item
		switch j {
		case '0':
			e = "⬛"
		case '1':
			e = "🔵"
		case '2':
			e = "🔴"
		case '3':
			e = "💙"
		case '4':
			e = "❤️"
		case directionSlice[0]: // Northwest
			e = movesSlice[0]
		case directionSlice[1]: // Northeast
			e = movesSlice[1]
		case directionSlice[2]: // Southwest
			e = movesSlice[2]
		case directionSlice[3]: // Southeast
			e = movesSlice[3]
		}

		// If the row is even
		if row%2 == 0 {
			formatted = append(formatted, "⬜"+e)
		} else { // If the row is odd
			formatted = append(formatted, e+"⬜")
		}

		// If we are at the end of the row
		if (i+1)%4 == 0 {
			formatted = append(formatted, "\n")
		}
	}

	// Add the bottom row of numbers
	formatted = append(formatted, "⏺1️⃣1️⃣2️⃣2️⃣3️⃣3️⃣4️⃣4️⃣")
	return strings.Join(formatted, "")
}

// Formats the user in a readable format
func formatUser(u *discordgo.User) string {
	return u.Username + "#" + u.Discriminator
}

// Generic message format for errors
func errorMessage(title string, message string) string {
	return "❌  **" + title + "**\n" + message
}

// Generic message format for successful operations
func successMessage(title string, message string) string {
	return "✅  **" + title + "**\n" + message
}

// Creates an embed for the game
func gameEmbed(s *discordgo.Session, cmd string, opponentID string, game *logic.Game, board string, spectate bool) *discordgo.MessageEmbed {
	opponent, err := s.User(opponentID)
	if err != nil {
		return &discordgo.MessageEmbed{
			Color:       c_RED,
			Description: "Error getting opponent",
		}
	}

	// Regular values
	color := c_BLUE
	status := "Your move"
	help := "For help type `!checkers help " + cmd + "`"
	cmdAndArgs := cmd + ":" + StringifyGame(opponentID, game)
	if spectate {
		// Spectator mode values
		color = c_DEFAULT
		status = "Waiting for opponent..."
		help = "For help type `!checkers help`"
		cmdAndArgs = "" // Ensures reactions won't do anything on old messages
	}

	// Shows the captured pieces
	p1score, p2score := logic.GetScore(game)
	var capturedPieces1 []string
	var capturedPieces2 []string
	for i := 0; i < p1score; i++ {
		capturedPieces1 = append(capturedPieces1, "🔴")
	}
	for i := 0; i < p2score; i++ {
		capturedPieces2 = append(capturedPieces2, "🔵")
	}

	// If there are no captured pieces set it to none so the embed is valid
	if len(capturedPieces1) == 0 {
		capturedPieces1 = []string{"None"}

	}
	if len(capturedPieces2) == 0 {
		capturedPieces2 = []string{"None"}
	}

	return &discordgo.MessageEmbed{
		Color: color,
		Title: "Checkers game against " + formatUser(opponent),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Status",
				Value: status,
			},
			{
				Name:  "Captured Red pieces",
				Value: strings.Join(capturedPieces1, ""),
			},
			{
				Name:  "Captured Blue pieces",
				Value: strings.Join(capturedPieces2, ""),
			},
			{
				Name:  "Board",
				Value: formatBoard(&board),
			},
			{
				Name:  "Help",
				Value: help,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: cmdAndArgs,
		},
	}
}
