package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// the de facto command handler (until i implement slash commands)
func MessageCreateEventHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Content[0] != '!' {
		return
	}
	cmd := strings.Split(e.Content[1:], " ")
	switch cmd[0] {
	case "ping":
		_, err := s.ChannelMessageSend(e.ChannelID, "pong!")
		if err != nil {
			log.Printf("failed to send pong: %v", err)
		}
	case "welcome":
		switch cmd[1] {
		case "channel":
			if cmd[2][0] == '<' {
				cmd[2] = cmd[2][2 : len(cmd[2])-1]
			}
			c, err := s.Channel(cmd[2])
			if err != nil {
				log.Printf("failed to fetch channel to set as welcome channel: %v", err)
				return
			}
			a.SetGuildWelcomeChannel(e.GuildID, c.ID)
		case "simulate":
			g, err := s.State.Guild(e.GuildID)
			if err != nil {
				log.Printf("failed to get guild from cache for welcome sim: %v", err)
				g, err = s.Guild(e.GuildID)
				if err != nil {
					log.Printf("failed to get guild from direct request for welcome sim: %v", err)
					return
				}
			}
			gw := a.GetGuildWelcome(g.ID)
			if gw.ChannelID != "" {
				wi := welcomeMessageInfo{
					"<@" + e.Author.ID + ">",
					e.Author.ID,
					e.Author.Username,
					g.Name,
					e.Author.AvatarURL(fmt.Sprint(PfpSize)),
				}
				welcome := generateWelcomeMessage(gw, wi)
				_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
				if err != nil {
					log.Printf("failed to send welcome sim: %v", err)
				}
			} else {
				s.ChannelMessageSend(e.ChannelID, "use `!welcome channel #channel` to set the welcome channel!")
			}
		default:
			_, err := s.ChannelMessageSend(e.ChannelID, "unknown welcome command!")
			if err != nil {
				log.Printf("failed to send message indicating unknown command: %v", err)
			}
		}
	}
}
