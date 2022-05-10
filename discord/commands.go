package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/logger"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "simple command to test if the bot is online",
		},
		{
			Name:        "welcome",
			Description: "commands related to welcome messages",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "set",
					Description: "set welcome message config; text placeholders: %guild%, %mention%, %username%, and %nickname%",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "the channel for the welcome message",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    false,
						},
						{
							Name:        "text",
							Description: "the main message text for the welcome message",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "type",
							Description: "the type of message (image or plain text) for the welcome message",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "image",
									Value: "image",
								},
								{
									Name:  "text",
									Value: "text",
								},
							},
							Required: false,
						},
						{
							Name:        "image",
							Description: "the background image for the welcome message",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "original",
									Value: "original",
								},
								{
									Name:  "beach",
									Value: "beach",
								},
								{
									Name:  "sleepy",
									Value: "sleepy",
								},
								{
									Name:  "friends",
									Value: "friends",
								},
								{
									Name:  "melon",
									Value: "melon",
								},
								{
									Name:  "sky",
									Value: "sky",
								},
							},
							Required: false,
						},
						{
							Name:        "imagetext",
							Description: "the text on the image for the welcome message",
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
				{
					Name:        "reset",
					Description: "reset welcome message config",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}

	commandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong!",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				logger.L.Error().Err(err).Msg("Failed to send interaction response")
			}
		},
		"welcome": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content strings.Builder
			g, err := s.Guild(i.GuildID)
			if err != nil {
				logger.L.Error().Err(err).Msg("Failed to get guild from direct request")
				return
			}
			if i.Interaction.Member.Permissions&discordgo.PermissionManageServer == discordgo.PermissionManageServer {
				switch i.ApplicationCommandData().Options[0].Name {
				case "set":
					content.WriteString("attempted to set: ")
					var attempt string
					if len(i.ApplicationCommandData().Options[0].Options) != 0 {
						for _, o := range i.ApplicationCommandData().Options[0].Options {
							switch o.Name {
							case "channel":
								cid := o.Value.(string)
								c, err := s.State.Channel(cid)
								if err != nil {
									logger.L.Warn().Err(err).Msg("Failed to get channel from state")
									c, err = s.Channel(cid)
									if err != nil {
										logger.L.Error().Err(err).Msg("Failed to get channel from direct request")
										return
									}
									s.State.ChannelAdd(c)
								}

								if c.Type != discordgo.ChannelTypeGuildText {
									content.WriteString("invalid channel, ")
								} else {
									database.SetGuildWelcomeChannel(i.GuildID, c.ID)
									content.WriteString("channel, ")
								}
							case "type":
								database.SetGuildWelcomeType(i.GuildID, o.StringValue())
								content.WriteString("type, ")
							case "text":
								database.SetGuildWelcomeText(i.GuildID, o.StringValue())
								content.WriteString("text, ")
							case "image":
								database.SetGuildWelcomeImage(i.GuildID, o.StringValue())
								content.WriteString("image, ")
							case "imagetext":
								database.SetGuildWelcomeImageText(i.GuildID, o.StringValue())
								content.WriteString("imagetext, ")
							}
						}
						attempt = content.String()
						attempt = attempt[:len(attempt)-2]
					} else {
						content.WriteString("nothing????")
						attempt = content.String()
					}
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: attempt,
						},
					})
					if err != nil {
						logger.L.Error().Err(err).Msg("Failed to send interaction response")
					}

				case "reset":
					content.WriteString("are you sure you want to reset your server's welcome config?")
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: content.String(),
							Components: []discordgo.MessageComponent{
								discordgo.ActionsRow{
									Components: []discordgo.MessageComponent{
										discordgo.Button{
											Emoji: discordgo.ComponentEmoji{
												Name: "💥",
											},
											CustomID: "reset_welcome",
											Label:    "reset",
											Style:    discordgo.DangerButton,
										},
									},
								},
							},
							Flags: 1 << 6,
						},
					})
					if err != nil {
						logger.L.Error().Err(err).Msg("Failed to send interaction response")
					}
				case "simu":
					u, err := s.GuildMember(g.ID, i.Member.User.ID)
					if err != nil {
						logger.L.Error().Err(err).Msg("Failed to get user from direct request")
						return
					}
					gm, err := s.GuildMembers(g.ID, "", 1000) // TODO: implement a function to get all members and utilize cache
					if err != nil {
						logger.L.Error().Err(err).Msgf("Failed to get guild members from direct request")
						return
					}

					gw := database.GetGuildWelcome(g.ID)
					if gw.ChannelID != "" {
						wi := welcomeMessageInfo{
							mention:   u.Mention(),
							nickname:  u.User.Username,
							username:  u.User.String(),
							guildName: g.Name,
							avatarURL: u.AvatarURL(fmt.Sprint(PfpSize)),
							members:   len(gm),
						}
						welcome := generateWelcomeMessage(gw, wi)
						_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
						if err != nil {
							logger.L.Error().Err(err).Msg("Failed to send welcome simulation")
						}
						content.WriteString("attempted to simulate welcome!")
					} else {
						content.WriteString("use `/welcome set channel` to set the welcome channel!")
					}
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: content.String(),
						},
					})
					if err != nil {
						logger.L.Error().Err(err).Msg("Failed to send interaction response")
					}
				}
			} else {
				content.WriteString("you do not have permission to use that command!")
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content.String(),
					},
				})
				if err != nil {
					logger.L.Error().Err(err).Msg("Failed to send interaction response")
				}
			}
		},
	}

	componentHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"reset_welcome": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Interaction.Member.Permissions&discordgo.PermissionManageServer == discordgo.PermissionManageServer {
				database.ResetGuild(i.GuildID)
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "guild welcome config reset!",
					},
				})
				if err != nil {
					logger.L.Error().Err(err).Msg("Failed to respond to interaction confirming database reset")
				}
			}
		},
	}
)
