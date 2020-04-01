package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jmsheff/discord-checkers/discord"
)

func main() {
	// Register the bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Missing token environment variable")
	}

	b, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err.Error())
	}

	// Register handlers
	b.AddHandler(discord.CommandsHandler)
	b.AddHandler(discord.ReactionsHandler)

	// Open a websocket connection to Discord and begin listening.
	err = b.Open()
	if err != nil {
		log.Panic("Could not connect to discord", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Print("Discord bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	b.Close()
}
