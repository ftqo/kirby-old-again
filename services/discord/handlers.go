package discord

import (
	"fmt"
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

/*

* figure out how to add a connection pool (will need to call from here)
* use same method to keep hoarded images
* finish welcome message generation
* benchmark if faster than java version, stop if not

 */

func ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("Bot is up!")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot joins a guild

}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild

}

func GuildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	// info := database.GuildWelcome{}
	// message, channel := welcome.GenerateWelcome(info, e.User)
	// _, err := s.ChannelMessageSendComplex(channel, &message)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(reflect.ValueOf(e))
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
