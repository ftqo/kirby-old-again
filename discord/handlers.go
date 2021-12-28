package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("Bot connected!")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot joins a guild
	a.InitServer(e.Guild.ID)
}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild
	a.CutServer(e.Guild.ID)
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
