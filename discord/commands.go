package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "the best ping command ever",
		},
		{
			Name:        "welcome",
			Description: "welcome",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Name:        "set",
					Description: "set welcome config",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "the channel u wanna set the welcome message to generate in",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    false,
						}, {
							Name:        "text",
							Description: "the channel u wanna set the welcome message to generate in",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						}, {
							Name:        "image",
							Description: "the channel u wanna set the welcome message to generate in",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "original",
									Value: "original",
								},
								{
									Name:  "grey",
									Value: "grey",
								},
							},
							Required: false,
						}, {
							Name:        "imagetext",
							Description: "the channel u wanna set the welcome message to generate in",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},

				{
					Name:        "simu",
					Description: "simulate a welcome message",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong!",
					Flags:   1 << 6,
				},
			})
		},
		"welcome": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content strings.Builder

			g, err := s.State.Guild(i.GuildID)
			if err != nil {
				log.Printf("failed to get guild from cache for welcome sim: %v", err)
				g, err = s.Guild(i.GuildID)
				if err != nil {
					log.Printf("failed to get guild from direct request for welcome sim: %v", err)
				}
			}
			if g.Permissions&discordgo.PermissionManageServer != discordgo.PermissionManageServer {
				switch i.ApplicationCommandData().Options[0].Name {
				case "set":
					content.WriteString("attempted to set")
					for _, o := range i.ApplicationCommandData().Options[0].Options {
						switch o.Name {
						case "channel":
							cid := o.Value.(string)
							c, err := s.State.Channel(cid)
							if err != nil {
								log.Printf("failed to get channel from cache: %v", err)
								c, err = s.Channel(cid)
								if err != nil {
									log.Printf("failed to get channel from direct request: %v", err)
								}
							}

							if c.Type != discordgo.ChannelTypeGuildText {
								content.WriteString(" - (bad) channel")
							} else {
								a.SetGuildWelcomeChannel(i.GuildID, c.ID)
								content.WriteString(" - channel")
							}
						case "text":
							a.SetGuildWelcomeText(i.GuildID, o.StringValue())
							content.WriteString(" - text")
						case "image":
							a.SetGuildWelcomeImage(i.GuildID, o.StringValue())
							content.WriteString(" - image")
						case "imagetext":
							a.SetGuildWelcomeImageText(i.GuildID, o.StringValue())
							content.WriteString(" - imagetext")
						}
					}

				case "simu":
					u, err := s.User(i.Member.User.ID)
					if err != nil {
						log.Printf("failed to get user from direct request for welcome simulation: %v", err)
					}
					gw := a.GetGuildWelcome(g.ID)
					if gw.ChannelID != "" {
						wi := welcomeMessageInfo{
							mention:   u.Mention(),
							nickname:  u.ID,
							username:  u.String(),
							guildName: g.Name,
							avatarURL: u.AvatarURL(fmt.Sprint(PfpSize)),
							members:   g.MemberCount,
						}
						welcome := generateWelcomeMessage(gw, wi)
						_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
						if err != nil {
							log.Printf("failed to send welcome simulation: %v", err)
						}
						content.WriteString("attempted to simulate welcome")
					} else {
						content.WriteString("use `/welcome set channel` to set the welcome channel!")
					}
				}
			} else {
				content.WriteString("you do not have permissions!")
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content.String(),
				},
			})
		},
	}
)
