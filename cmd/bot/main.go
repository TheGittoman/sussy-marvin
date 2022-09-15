package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheGittoman/sussy-marvin/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	BotID string
	AppID string
)

func init() {
	const fileName = "./config/config.json"

	cfg, err := config.ParseConfigFromJSONFile(fileName)
	if err != nil {
		panic(err)
	}
	flag.StringVar(&Token, "t", cfg.Token, "Bot Token")
	flag.StringVar(&AppID, "t", cfg.AppID, "Application ID")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := "<@" + s.State.User.ID + "> "
	fmt.Print("Message sent! GuildID: " + m.GuildID + "\n")

	if m.Author.ID == s.State.User.ID {
		fmt.Print("User and Marvin are the same person\n")
		return
	}

	if m.Content == prefix+"pong" {
		s.ChannelMessageSend(m.ChannelID, "PING!")
	}
	if m.Content == prefix+"ping" {
		s.ChannelMessageSend(m.ChannelID, "PONG!")
	}
}
