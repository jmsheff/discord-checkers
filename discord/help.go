package discord

import "github.com/bwmarrin/discordgo"

func helpCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, topic string) {
	var title string
	var description string
	var fields []*discordgo.MessageEmbedField

	switch topic {
	case "invites":
		title = "✉️  Invites - Checkers Help"
		description = "Invites allow you to start a game with a player. Invites CANNOT be sent:\n  • Through DM\n  • By or to bots\n  • To yourself\nSee below for available commands."
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  "General invites",
				Value: "`!checkers invite`: Sends an invite to the channel the command was sent in for any user to accept",
			},
			{
				Name:  "Direct invites",
				Value: "`!checkers invite @<user>`: Sends an invite directly to the mentioned user.",
			},
		}
	case "select":
		title = "⏺  Selection - Checkers Help"
		description = "Selection is done through reactions. Valid selection reactions are set by the bot. Select a piece by reacting with the corresponding emojis."
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Example",
				Value: "To select the piece at the square F1 react with  🇫  followed by  1️⃣ .",
			},
			{
				Name:  "Confirmation",
				Value: "Confirm a selection by reacting with a  ✅ . This will validate your selection and show the available moves for the selected piece on the board.",
			},
		}
	case "move":
		title = "↗️  Movement - Checkers Help"
		description = "Movement is done through reactions. Valid movement reactions are set by the bot. The moves for the piece you selected are shown on the board. Move the piece by reacting with the corresponding emoji."
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Example",
				Value: "To move the selected piece to the square Northeast of itself, react with a  ↗️",
			},
			{
				Name:  "Cancel",
				Value: "To select a different piece, react with a  ❌ . This will bring you back to the selection step.",
			},
		}
	default:
		title = "ℹ️  Topics - Checkers Help"
		description = "Pick a topic below to get help"
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  "✉️  Invites",
				Value: "`!checkers help invites`:  Gives instruction on how to send invites ",
			},
			{
				Name:  "⏺  Selection",
				Value: "`!checkers help select`: Gives help on how to select a piece",
			},
			{
				Name:  "↗️  Movement",
				Value: "`!checkers help move`: Provides help on how to move a piece",
			},
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Fields:      fields,
		Color:       c_BLUE,
	})
}
