package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// the de facto command handler (until i implement slash commands)
// update: this is starting to look like a shitshow...
func MessageCreateEventHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Content == "" {
		return
	}
	if e.Content[0] != '!' {
		return
	}

	cmd := strings.Fields(e.Content[1:])
	switch cmd[0] {
	case "ping":
		_, err := s.ChannelMessageSend(e.ChannelID, "pong!")
		if err != nil {
			log.Printf("failed to send pong: %v", err)
		}
	case "welcome":
		_, err := s.GuildMember(e.GuildID, e.Author.ID)
		if err != nil {
			log.Printf("failed to fetch guild member for permissions check: %v", err)
		}
		mPerms, err := s.State.MessagePermissions(e.Message)
		if err != nil {
			log.Printf("failed to get message permissions for guild member in state: %v", err)
		}
		if mPerms&discordgo.PermissionManageServer != discordgo.PermissionManageServer {
			return
		}
		g, err := s.State.Guild(e.GuildID)
		if err != nil {
			log.Printf("failed to get guild from cache for welcome sim: %v", err)
			g, err = s.Guild(e.GuildID)
			if err != nil {
				log.Printf("failed to get guild from direct request for welcome sim: %v", err)
				return
			}
		}
		switch cmd[1] {
		case "set":
			switch cmd[2] {
			case "channel":
				if cmd[3][0] == '<' {
					cmd[3] = cmd[3][2 : len(cmd[3])-1]
				}
				c, err := s.State.Channel(cmd[3])
				if err != nil {
					log.Printf("failed to get channel from cache for welcome set channel: %v", err)
					c, err = s.Channel(cmd[3])
					if err != nil {
						log.Printf("failed to get channel from direct request: %v", err)
						return
					}
				}
				a.SetGuildWelcomeChannel(e.GuildID, c.ID)
			case "text":
				beg := strings.Index(e.Content, "\"")
				end := strings.LastIndex(e.Content, "\"")
				if beg == -1 || end == -1 || beg >= end {
					_, err := s.ChannelMessageSend(e.ChannelID, "put **straight double quotes** around the text you want to set as your welcome text!")
					if err != nil {
						log.Printf("failed to send message demanding the user add quotes to their welcome text: %v", err)
					}
					return
				}
				text := e.Content[beg+1 : end]
				a.SetGuildWelcomeText(e.GuildID, text)
			case "image":
				a.SetGuildWelcomeImage(e.GuildID, cmd[3])
			case "imagetext":
				beg := strings.Index(e.Content, "\"")
				end := strings.LastIndex(e.Content, "\"")
				if beg == -1 || end == -1 || beg >= end {
					_, err := s.ChannelMessageSend(e.ChannelID, "put **straight double quotes** around the text you want to set as your welcome image text!")
					if err != nil {
						log.Printf("failed to send message demanding the user add quotes to their welcome image text: %v", err)
					}
					return
				}
				imagetext := e.Content[beg+1 : end]
				a.SetGuildWelcomeImageText(e.GuildID, imagetext)
			}
			_, err := s.ChannelMessageSend(e.ChannelID, "guild welcome setting updated (hopefully)!")
			if err != nil {
				log.Printf("failed to send message confirming updated guild welcome settings: %v", err)
			}
		case "simu":
			gw := a.GetGuildWelcome(g.ID)
			if gw.ChannelID != "" {
				wi := welcomeMessageInfo{
					mention:   e.Author.Mention(),
					nickname:  e.Author.ID,
					username:  e.Author.String(),
					guildName: g.Name,
					avatarURL: e.Author.AvatarURL(fmt.Sprint(PfpSize)),
					members:   g.MemberCount,
				}
				welcome := generateWelcomeMessage(gw, wi)
				_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
				if err != nil {
					log.Printf("failed to send welcome sim: %v", err)
				}
			} else {
				_, err := s.ChannelMessageSend(e.ChannelID, "use `!welcome set channel #channel` to set the welcome channel!")
				if err != nil {
					log.Printf("failed to send message demanding the user use the correct command: %v", err)
				}
			}
		default:
			_, err := s.ChannelMessageSend(e.ChannelID, "unknown welcome command!")
			if err != nil {
				log.Printf("failed to send message indicating unknown command: %v", err)
			}
		}
	}
}
