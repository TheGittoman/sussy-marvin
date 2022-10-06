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

func checkError(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

func removeInteraction(s *discordgo.Session, i *discordgo.Interaction) {
	time.Sleep(time.Second * 5)
	err := s.InteractionResponseDelete(i)
	checkError(err)
}

// read a random line including marvin and returns it
func readData(find string, filename string) (outString string) {
	rand.Seed(time.Now().Unix())
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
		{ // my-command
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
		{ // quote
			Name:        "quote",
			Description: "get random quote from marvin the robot",
		},
		{ // delete-messages
			Name:        "delete-messages",
			Description: "delete messages up to an amount",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "delete-up-to",
					Description: "delete up to amount of comments",
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
		"quote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content := readData("Marvin", "./internal/data/data.dat")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"delete-messages": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			permission := i.Interaction.Member.Permissions // get interaction user perms

			if permission&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages { // if permissions match manage message bit flags
				options := i.ApplicationCommandData().Options

				if int(options[0].IntValue()) > 100 || int(options[0].IntValue()) < 0 { // checks that the option range is correct, must be some more elegant way of doing this
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Invalid input: number must positive and less or equal to 100"),
						},
					})
				} else {
					messages, err := s.ChannelMessages(i.ChannelID, int(options[0].IntValue()), "", "", "")
					checkError(err)
					var messagesToString []string
					for _, v := range messages {
						messagesToString = append(messagesToString, v.ID)
					}
					err = s.ChannelMessagesBulkDelete(i.ChannelID, messagesToString)
					checkError(err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("messages deleted: %d", int(options[0].IntValue())),
						},
					})
				}
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Not permitted to run this command"),
					},
				})
			}
			removeInteraction(s, i.Interaction)
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
