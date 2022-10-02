package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheGittoman/sussy-marvin/internal/atCommands"
	"github.com/TheGittoman/sussy-marvin/internal/config"
	"github.com/TheGittoman/sussy-marvin/internal/slashCommands"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

// Variables used for command line parameters
func init() {
	fmt.Print("Bot Initialized!\n")
}

func main() {
	// Create a new Discord session using the provided bot token.

	slashCommands.Test()
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
	dg.AddHandler(atCommands.MessageCreate)

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
