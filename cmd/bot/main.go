package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/TheGittoman/sussy-marvin/internal/atCommands"
	"github.com/TheGittoman/sussy-marvin/internal/config"
	"github.com/TheGittoman/sussy-marvin/internal/slashCommands"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session
var err error
var cfg config.Config

// Variables used for command line parameters
func init() {
	const fileName = "config/config.json"
	cfg, err := config.ParseConfigFromJSONFile(fileName)
	if err != nil {
		fmt.Println("Error reading the config.json, ", err)
		return
	}
	s, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := slashCommands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommands.Commands))
	for i, v := range slashCommands.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, cfg.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(atCommands.MessageCreate)

	// stops the bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *&cfg.RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, cfg.GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")

}
