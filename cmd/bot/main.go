package main

import (
	"github.com/TheGittoman/sussy-marvin/internal/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	const fileName = "./config/config.json"

	cfg, err := config.ParseConfigFromJSONFile(fileName)
	if err != nil {
		panic(err)
	}

	s, err := discordgo.New("Bot" + cfg.Token)
	if err != nil {
		panic(err)
	}

	if err = s.Open(); err != nil {
		panic(err)
	}
}
