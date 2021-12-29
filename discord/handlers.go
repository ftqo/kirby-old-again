package discord

import (
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
		log.Panicln(err)
	}
	log.Println("Bot connected!")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot turns on or joins a guild
	a.InitServer(e.Guild.ID)
}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild
	if !e.Unavailable {
		a.CutServer(e.Guild.ID)
	}
}

func GuildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	// if GuildWelcome is cached, use it
	// else
	gw := a.GetGuildWelcome(e.GuildID)
	msg := GenerateWelcome(gw, e.User)
	s.ChannelMessageSendComplex(gw.Channel, &msg)
}

func MessageCreateEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

}
