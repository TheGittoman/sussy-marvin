package atCommands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GetMessage(m *discordgo.MessageCreate, s *discordgo.Session, find string, message string) {
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
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	GetMessage(m, s, "testi", "Testi vastaus")
	GetMessage(m, s, "ping", "PONG!")
	GetMessage(m, s, "pong", "PING!")
}
