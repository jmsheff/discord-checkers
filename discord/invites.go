package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jmsheff/discord-checkers/logic"
)

// Handlers/Functions for everything invite related

// Sends a invite to game to a users DM
func sendDirectInvite(s *discordgo.Session, m *discordgo.MessageCreate, recipient *discordgo.User) {
	if m.Author.ID == recipient.ID {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Invalid recipient", "Cannot play against yourself!"))
		return
	}

	if recipient.Bot {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Invalid recipient", "Cannot play against bot!"))
		return
	}

	dm, err := s.UserChannelCreate(recipient.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Bot error", "Error creating direct message."))
		return
	}

	invite, err := s.ChannelMessageSendEmbed(dm.ID, &discordgo.MessageEmbed{
		Title:       "Checkers game invite from " + formatUser(m.Author),
		Description: "Click the  ✅  to accept this invitation, or the  ❌  to deny.",
		Color:       c_BLUE,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "invite:" + m.Author.ID,
		},
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Bot error", "Error sending invite."))
		return
	}

	s.MessageReactionAdd(dm.ID, invite.ID, "✅")
	s.MessageReactionAdd(dm.ID, invite.ID, "❌")

	s.ChannelMessageSend(m.ChannelID, successMessage("Success", "Invite sent to "+formatUser(recipient)+"!"))
}

// Sends a general invite for any user in the channel to accept
func sendGeneralInvite(s *discordgo.Session, m *discordgo.MessageCreate) {
	invite, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Checkers game invite from " + formatUser(m.Author),
		Description: "Click the  ✅  to accept this invitation.",
		Color:       c_BLUE,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "generalinvite:" + m.Author.ID,
		},
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Bot error", "Error sending invite."))
		return
	}

	s.MessageReactionAdd(m.ChannelID, invite.ID, "✅")
}

// Handles all invite related commands
func inviteCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Bot error", "Error getting channel."))
		return
	}

	// Ensure that the command is not being sent from a dm
	if c.Type == discordgo.ChannelTypeDM {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Invalid channel", "Cannot send invites from a DM"))
		return
	}

	recipients := m.Mentions
	if len(recipients) == 1 {
		sendDirectInvite(s, m, recipients[0])
	} else if len(recipients) == 0 {
		// Ensure this is not a mistake by making sure these are the only 2 arguments
		if len(cmd) == 1 {
			sendGeneralInvite(s, m)
		} else {
			s.ChannelMessageSend(m.ChannelID, errorMessage("Invalid Reciepient", "Ensure you are mentioning the player in the format of @<user>. Or, if you are trying to send a general invite leave the user blank."))
		}
	} else if len(recipients) > 1 {
		s.ChannelMessageSend(m.ChannelID, errorMessage("Invalid invite", "Cannot invite multiple players!"))
	}
}

// Handles all invite related reactions
func inviteReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd, m *discordgo.Message, user *discordgo.User, opponentID string, general bool) {
	// If the reaction comes from the sender of the invite(This will only happen in the case of general invites)
	if r.UserID == opponentID {
		return
	}
	sender, err := s.User(opponentID)
	if err != nil || sender == nil {
		return
	}
	opponentDM, _ := s.UserChannelCreate(opponentID)
	if r.Emoji.Name == "✅" && (general || !hasOtherReactionsBesides("✅", m.Reactions)) {
		s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &discordgo.MessageEmbed{
			Title:       "Invite Accepted!",
			Description: "Invite from " + formatUser(sender) + " accepted!",
			Color:       c_GREEN,
		})

		// Create a game object
		game := logic.Game{
			Selected: 0,
			Board:    "11111111111100000000222222222222",
			Turn:     2,
		}

		var reciepientDMID string
		if !general {
			reciepientDMID = r.ChannelID
		} else {
			reciepientDM, _ := s.UserChannelCreate(r.UserID)
			reciepientDMID = reciepientDM.ID
		}

		gamemsg, err := s.ChannelMessageSendEmbed(reciepientDMID, gameEmbed(s, "select", opponentID, &game, game.Board, false))
		if err != nil {
			return
		}
		s.ChannelMessageSend(opponentDM.ID, successMessage("Game on!", formatUser(user)+" accepted your checkers invite! Wait here for them to make their move."))
		addSelectReactions(s, reciepientDMID, gamemsg.ID, &game)
	} else if !general && r.Emoji.Name == "❌" && !hasOtherReactionsBesides("❌", m.Reactions) {
		s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &discordgo.MessageEmbed{
			Title:       "Invite Declined",
			Description: "Invite from " + formatUser(sender) + " declined.",
			Color:       c_RED,
		})
		s.ChannelMessageSend(opponentDM.ID, errorMessage("Invite declined", formatUser(user)+" declined your checkers game invite."))
	}
}
