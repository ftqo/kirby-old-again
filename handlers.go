package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func readyHandler(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("Bot is up!")
}

func guildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot joins a guild

}

func guildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild

}

func guildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {

}

func guildMemberRemoveEventHandler(s *discordgo.Session, e *discordgo.GuildMemberRemove) {

}

func messageDeleteEventHandler(s *discordgo.Session, e *discordgo.MessageDelete) {

}

func messageReactionAddEventHandler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {

}

func messageReactionRemoveEventHandler(s *discordgo.Session, e *discordgo.MessageReactionRemove) {

}

func messageCreateEventHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Content == "!ping" {
		_, err := s.ChannelMessageSend(e.ChannelID, "pong!")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func channelDeleteEventHandler(s *discordgo.Session, e *discordgo.ChannelDelete) {

}
