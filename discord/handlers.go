package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	usd := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "the stars",
				Type: 3,
			},
		},
	}
	err := s.UpdateStatusComplex(usd)
	if err != nil {
		log.Panicf("failed to update status: %v", err)
	}
	log.Print("connected to discord !")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot turns on or joins a guild
	a.InitGuild(e.Guild.ID)
}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild
	if !e.Unavailable {
		a.CutGuild(e.Guild.ID)
	}
}

func GuildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	g, err := s.State.Guild(e.GuildID)
	if err != nil {
		log.Printf("failed to get guild from cache when GuildMemberAdd was fired: %v", err)
		g, err = s.Guild(e.GuildID)
		if err != nil {
			log.Printf("failed to get guild from direct request: %v", err)
			return
		}
	}
	gw := a.GetGuildWelcome(g.ID)
	if gw.ChannelID != "" {
		wi := welcomeMessageInfo{
			mention:   e.User.Mention(),
			nickname:  e.User.Username,
			username:  e.User.String(),
			guildName: g.Name,
			avatarURL: e.User.AvatarURL(fmt.Sprint(PfpSize)),
			members:   g.MemberCount,
		}
		welcome := generateWelcomeMessage(gw, wi)
		_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
		if err != nil {
			log.Printf("failed to send welcome message: %v", err)
		}
	}
}

func ChannelDeleteEventHandler(s *discordgo.Session, e *discordgo.ChannelDelete) {
	gw := a.GetGuildWelcome(e.GuildID)
	if e.Channel.ID == gw.ChannelID {
		a.SetGuildWelcomeChannel(e.GuildID, "")
	}
}

func InteractionCreateEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	case discordgo.InteractionMessageComponent:

		if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
			h(s, i)
		}
	}
}
