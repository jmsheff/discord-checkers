package discord

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jmsheff/discord-checkers/logic"
)

// Color Enum
const (
	c_DEFAULT             = 0
	c_AQUA                = 1752220
	c_GREEN               = 3066993
	c_BLUE                = 3447003
	c_PURPLE              = 10181046
	c_GOLD                = 15844367
	c_ORANGE              = 15105570
	c_RED                 = 15158332
	c_GREY                = 9807270
	c_DARKER_GREY         = 8359053
	c_NAVY                = 3426654
	c_DARK_AQUA           = 1146986
	c_DARK_GREEN          = 2067276
	c_DARK_BLUE           = 2123412
	c_DARK_PURPLE         = 7419530
	c_DARK_GOLD           = 12745742
	c_DARK_ORANGE         = 11027200
	c_DARK_RED            = 10038562
	c_DARK_GREY           = 9936031
	c_LIGHT_GREY          = 12370112
	c_DARK_NAVY           = 2899536
	c_LUMINOUS_VIVID_PINK = 16580705
	c_DARK_VIVID_PINK     = 12320855
)

// Get the reactions preset by the bot
func getBotReactions(reactions []*discordgo.MessageReactions) []*discordgo.MessageReactions {
	var botReactions []*discordgo.MessageReactions
	// Filters reactions to only bot reactions
	for _, r := range reactions {
		if r.Me {
			botReactions = append(botReactions, r)
		}
	}

	return botReactions
}

// Check if users reaction is one preset by the bot
func isBotReaction(s *discordgo.Session, reactions []*discordgo.MessageReactions, emoji *discordgo.Emoji) bool {
	for _, r := range reactions {
		if r.Emoji.Name == emoji.Name && r.Me {
			return true
		}
	}

	return false
}

// Good for making sure only 1 reaction is selected to prevent reaction spam
func hasOtherReactionsBesides(allowed string, reactions []*discordgo.MessageReactions) bool {
	for _, r := range reactions {
		// If there is more than one reaction on a reaction that isn't the one allowed
		if r.Count > 1 && allowed != r.Emoji.Name {
			// If it has a count greater than 1
			return true
		}
	}

	return false
}

// Parses a game string and returns the opponent ID and the Game object
func ParseGame(s string) (string, logic.Game, error) {
	var game logic.Game
	values := strings.Split(s, " ")

	if len(values) != 4 {
		return "", logic.Game{}, errors.New("Invalid input")
	}

	turn, err := strconv.ParseUint(values[1], 10, 0)
	if err != nil {
		return "", logic.Game{}, errors.New("Could not parse turn")
	}

	game.Turn = uint8(turn)
	game.Board = values[2]
	selected, err := strconv.ParseUint(values[3], 10, 0)
	if err != nil {
		return "", logic.Game{}, errors.New("Could not parse selected")
	}
	game.Selected = uint8(selected)

	return values[0], game, nil
}

// Stringifies a game
func StringifyGame(opponentID string, game *logic.Game) string {
	values := []string{
		opponentID,
		strconv.FormatUint(uint64(game.Turn), 10),
		game.Board,
		strconv.FormatUint(uint64(game.Selected), 10),
	}

	return strings.Join(values, " ")
}
