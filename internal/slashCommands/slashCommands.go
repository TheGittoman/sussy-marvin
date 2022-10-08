package slashCommands

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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
		{ // quote
			Name:        "quote",
			Description: "get random quote from marvin the robot",
		},
		{ // delete-message
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
		{ // tag-color
			Name:        "tagcolor",
			Description: "delete messages up to an amount",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "hex-color",
					Description: "hex code of role",
					Required:    true,
				},
			},
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
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
			// httpError := make(map[string]map[string]string)
			// httpError["HTTP 400 Bad Request"]["message"] = "You can only bulk delete messages that are under 14 days old."
			content := ""
			fail := false

			if permission&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages { // if permissions match manage message bit flags
				options := i.ApplicationCommandData().Options

				if int(options[0].IntValue()) > 100 || int(options[0].IntValue()) < 0 {
					content += "Invalid input: input < 100 or > 0\n"
					fail = true
				}

				if !fail { // if there is no fails continue with the command
					messages, err := s.ChannelMessages(i.ChannelID, int(options[0].IntValue()), "", "", "")
					if err != nil {
						content += err.Error()
					}

					var messagesToString []string
					for _, v := range messages {
						messagesToString = append(messagesToString, v.ID)
					}

					err = s.ChannelMessagesBulkDelete(i.ChannelID, messagesToString)
					if err != nil {
						for _, message := range messagesToString {
							s.ChannelMessageDelete(i.ChannelID, message)
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: fmt.Sprintf("%s \n deleting one by one %d messages", err.Error(), options[0].IntValue()),
								},
							})
						}
					}

					/* HTTP 400 Bad Request,
					{"message": "You can only bulk delete messages that are under 14 days old.",
					"code": 50034} this message will be deleted in 5 seconds
					*/

					if err == nil {
						content += "messages deleted: " + strconv.Itoa(int(options[0].IntValue()))
					}

				}
			} else {
				content += "you are not permitted to run this command"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content + " | this message will be deleted in 5 seconds"),
				},
			})
			removeInteraction(s, i.Interaction)
		},
		"tagcolor": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			option := i.ApplicationCommandData().Options
			hex_code_string := hex.EncodeToString([]byte(option[0].StringValue()))
			hex_code_int, err := strconv.ParseInt(hex_code_string, 8, 8)
			if err != nil {
				log.Println(err)
			}
			role := new(discordgo.Role)
			role.Name = "test"
			role.Color = int(hex_code_int)
			role.Mentionable = false
			role.Managed = false

			s.State.RoleAdd(i.GuildID, role)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Tag color changed to #%s", hex_code_string),
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

// {
// 	Name:        "basic-command",
// 	Description: "description",
// },
// {
// 	Name:        "basic-command-with-files",
// 	Description: "description",
// },
// {
// 	Name:        "followups",
// 	Description: "description",
// },
// {
// 	Name:        "localized-command",
// 	Description: "description",
// },
// {
// 	Name:        "options",
// 	Description: "description",
// },
// {
// 	Name:        "permission-overview",
// 	Description: "description",
// },
