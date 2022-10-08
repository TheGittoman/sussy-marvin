package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/TheGittoman/sussy-marvin/internal/atCommands"
	"github.com/TheGittoman/sussy-marvin/internal/config"
	"github.com/TheGittoman/sussy-marvin/internal/slashCommands"
	"github.com/bwmarrin/discordgo"
)

func botInitializing() (s_ *discordgo.Session, cfg_ *config.Config, err_ error) {
	const fileName = "config/config.json"
	cfg_, err_ = config.ParseConfigFromJSONFile(fileName)
	if err_ != nil {
		log.Println("Error reading the config.json, ", err_)
		return
	}

	s_, err_ = discordgo.New("Bot " + cfg_.Token)
	if err_ != nil {
		log.Println("error creating Discord session,", err_)
		return
	}
	s_.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := slashCommands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	return s_, cfg_, err_
}

func main() {
	var err error
	var cfg *config.Config
	s, cfg, err := botInitializing() // initialize variables needed for running the bot
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	err = s.Open() // open connection to discord web services and start listening
	if err != nil {
		log.Println("error opening connection,", err)
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

	if cfg.RemoveCommands {
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
	s.Close()
}
