package slashCommands

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// read a random line including marvin and returns it
func readData(find string) (outString string) {
	rand.Seed(time.Now().Unix())
	var filename string = "./internal/data/data.dat"
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	sc := bufio.NewScanner(file)
	lines := make([]string, 0)

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	var marvin []string

	for _, v := range lines {
		if strings.Contains(v, find) {
			marvin = append(marvin, v)
		}
	}

	outString = marvin[func() int {
		var i int
		i = rand.Intn(len(marvin))
		return i
	}()]

	return
}

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
		{
			Name:        "quote",
			Description: "get random quote from marvin the robot",
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
		"quote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content := readData("Marvin")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
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
