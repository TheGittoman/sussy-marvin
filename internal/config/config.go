package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token          string `json:"token"`
	GuildID        string `json:"GuildID"`
	ChannelID      string `json:"ChannelID"`
	AppID          string `json:"AppID"`
	Prefix         string `json:"prefix"`
	RemoveCommands bool   `json:"RemoveCommands"`
}

func ParseConfigFromJSONFile(fileName string) (c *Config, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}

	c = new(Config)
	err = json.NewDecoder(f).Decode(c)
	return
}
