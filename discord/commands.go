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
		case "set":
			switch cmd[2] {
			case "channel":
				if cmd[3][0] == '<' {
					cmd[3] = cmd[3][2 : len(cmd[3])-1]
				}
				c, err := s.Channel(cmd[3])
				if err != nil {
					log.Printf("failed to fetch channel to set as welcome channel: %v", err)
					return
				}
				a.SetGuildWelcomeChannel(e.GuildID, c.ID)
			case "text":
				text := e.Content[strings.Index(e.Content, "\" ")+len("\" ") : strings.LastIndex(e.Content, "\"")]
				a.SetGuildWelcomeText(e.GuildID, text)
			case "image":
				a.SetGuildWelcomeImage(e.GuildID, cmd[3])
			case "imagetext":
				imagetext := e.Content[strings.Index(e.Content, "\" ")+len("\" ") : strings.LastIndex(e.Content, "\"")]
				a.SetGuildWelcomeImageText(e.GuildID, imagetext)
			}
			_, err := s.ChannelMessageSend(e.ChannelID, "guild welcome setting updated (hopefully)!")
			if err != nil {
				log.Printf("failed to send message confirming updated guild welcome settings: %v", err)
			}
		case "simu":
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
					mention:   "<@" + e.Author.ID + ">",
					nickname:  e.Author.ID,
					username:  e.Author.Username + "#" + e.Author.Discriminator,
					guildName: g.Name,
					avatarURL: e.Author.AvatarURL(fmt.Sprint(PfpSize)),
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
