package slashCommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	IntegerOptionMinValue          = 1.0
	DmPermission                   = false
	DefaultMemberPermissions int64 = discordgo.PermissionManageServer

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "my-command",
			Description: "My first command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string",
					Description: "string to print",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"my-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			content := ""

			if option, ok := optionMap["string"]; ok {
				margs = append(margs, option.StringValue())
				content += "%s\n"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						content,
						margs...,
					),
				},
			})

		},
	}
)

// type Commands struct {
// 	IntegerOptionMinValue    json.Number     `json:"IntegerOptionMinValue`
// 	DmPermission             bool            `json:"DmPermission"`
// 	CefaultMemberPermissions int64           `json:"DefaultMemberPermissions`
// 	CommandsFrame            []CommandsFrame `json:"CommandsFrame"`
// }

// type CommandsFrame struct {
// 	Name        string `json:"Name"`
// 	Description string `json:"Description"`
// }

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func readCommands() (coms *Commands) {
// 	file, err := os.ReadFile("./config/commands.json")
// 	check(err)
// 	json.Unmarshal(file, &coms)
// 	if coms.DefaultMemberPermissions == 0 {
// 		coms.DefaultMemberPermissions = discordgo.PermissionManageServer
// 	}
// 	return coms
// }

// func Test() {
// 	var com = readCommands()
// 	println(com.DefaultMemberPermissions)
// }
