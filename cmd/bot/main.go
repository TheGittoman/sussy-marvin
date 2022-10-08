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

func botInitializing(s *discordgo.Session, err error, cfg *config.Config) (s_ *discordgo.Session, err_ error, cfg_ *config.Config) {
	const fileName = "config/config.json"
	cfg, err = config.ParseConfigFromJSONFile(fileName)
	log.Printf("%t", cfg.RemoveCommands)

	if err != nil {
		log.Println("Error reading the config.json, ", err)
		return
	}
	s, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := slashCommands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	return s, err, cfg
}

func init() {
}

var run bool = true

func main() {
	var s *discordgo.Session
	var err error
	var cfg *config.Config
	s, err, cfg = botInitializing(s, err, cfg) // initialize variables needed for running the bot

	if run {

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
	if !run {
		log.Println("Not running the bot")
	}
}
