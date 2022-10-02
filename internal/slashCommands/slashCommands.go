package slashCommands

import (
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	IntegerOptionMinValue    json.Number `json:"IntegerOptionMinValue`
	DmPermission             bool        `json:"DmPermission"`
	DefaultMemberPermissions int64
	CommandsFrame            []CommandsFrame `json:"CommandsFrame"`
}

type CommandsFrame struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readCommands() (coms *Commands) {
	file, err := os.ReadFile("./config/commands.json")
	check(err)
	json.Unmarshal(file, &coms)
	coms.DefaultMemberPermissions = discordgo.PermissionManageServer
	return coms
}

func Test() {
	var com = readCommands()
	println(com.DefaultMemberPermissions)
}
