package slashCommands

import (
	"bufio"
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

func parseHexString(hexString string) string {
	hexString = strings.Replace(hexString, "#", "", -1)
	return hexString
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
		{ // delete-mesage
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
			Description: "change tagcolor based on hex value",
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
			hexString := option[0].StringValue() // get the hex code
			boolFail := false
			hexInt, err := strconv.ParseInt(parseHexString(hexString), 16, 64)
			if err != nil {
				log.Println(err)
				boolFail = true
			}
			hexInt_ := int(hexInt)
			boolFalse := false

			if len(parseHexString(hexString)) < 6 || len(parseHexString(hexString)) > 6 {
				boolFail = true
			}

			if !boolFail {

				roleParams := new(discordgo.RoleParams)
				roleParams.Name = "color"
				roleParams.Color = &hexInt_
				roleParams.Hoist = &boolFalse
				roleParams.Mentionable = &boolFalse

				roles, err := s.GuildRoles(i.GuildID)
				if err != nil {
					log.Println(err)
				}

				memberRoles := i.Interaction.Member.Roles
				for _, r := range roles {
					for _, v := range memberRoles {
						if r.ID == v {
							err = s.GuildMemberRoleRemove(i.GuildID, i.Interaction.Member.User.ID, v)
							if err != nil {
								log.Println(err)
							}
						}
					}
				}

				foundRoleName := false
				for _, r := range roles {
					if r.Color == *roleParams.Color {
						foundRoleName = true
						s.GuildMemberRoleAdd(i.GuildID, i.Interaction.Member.User.ID, r.ID)
						break
					}
				}

				if !foundRoleName {
					role, err := s.GuildRoleCreate(i.GuildID, roleParams)
					if err != nil {
						log.Println(err)
					}
					s.GuildMemberRoleAdd(i.GuildID, i.Interaction.Member.User.ID, role.ID)
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Tag color of user %s changed to #%s", i.Interaction.Member.User.Username, hexString),
					},
				})
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Input is incorrect",
				},
			})
			removeInteraction(s, i.Interaction)
		},
	}
)
