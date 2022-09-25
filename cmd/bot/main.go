package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/TheGittoman/sussy-marvin/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
func init() {
	fmt.Print("Bot Initialized!\n")
}

func main() {
	// Create a new Discord session using the provided bot token.

	const fileName = "config/config.json"
	cfg, err := config.ParseConfigFromJSONFile(fileName)
	if err != nil {
		fmt.Println("Error reading the config.json, ", err)
		return
	}
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Sussy Marvin is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func getMessage(m *discordgo.MessageCreate, s *discordgo.Session, find string, message string) {
	if m.Content == "" {
		return
	}
	messageContent := m.Content
	if strings.Contains(messageContent, find) {
		s.ChannelMessageSend(m.ChannelID, message)
	}
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Print("Message sent! GuildID: " + m.GuildID + "\n")

	if m.Author.ID == s.State.User.ID {
		fmt.Print("User and Marvin are the same person\n")
		return
	}

	getMessage(m, s, "testi", "Testi vastaus")
	getMessage(m, s, "ping", "PONG!")
	getMessage(m, s, "pong", "PING!")
}
