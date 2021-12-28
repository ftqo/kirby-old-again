package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("Bot connected!")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot joins a guild

}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild

}

func GuildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	// pp, err := a.DB.Begin(context.Background())
	// pp.Prepare()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// info := database.GuildWelcome{}
	// message, channel := GenerateWelcome(info, e.User)
	// _, err = s.ChannelMessageSendComplex(channel, &message)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func GuildMemberRemoveEventHandler(s *discordgo.Session, e *discordgo.GuildMemberRemove) {

}

func MessageDeleteEventHandler(s *discordgo.Session, e *discordgo.MessageDelete) {

}

func MessageReactionAddEventHandler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {

}

func MessageReactionRemoveEventHandler(s *discordgo.Session, e *discordgo.MessageReactionRemove) {

}

func MessageCreateEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

}

func ChannelDeleteEventHandler(s *discordgo.Session, e *discordgo.ChannelDelete) {

}
