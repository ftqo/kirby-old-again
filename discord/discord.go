package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/fittsqo/kirby/database"
	"github.com/fittsqo/kirby/files"
)

var s *discordgo.Session
var f *files.Hoarder
var a *database.Adapter

func Start(token string) {
	var err error
	f = new(files.Hoarder)
	f.LoadImages()
	a = database.Connect()

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}
	s.AddHandler(ReadyHandler)
	s.AddHandler(GuildCreateEventHandler)
	s.AddHandler(GuildDeleteEventHandler)
	s.AddHandler(GuildMemberAddEventHandler)
	s.AddHandler(GuildMemberRemoveEventHandler)
	s.AddHandler(MessageDeleteEventHandler)
	s.AddHandler(MessageReactionAddEventHandler)
	s.AddHandler(MessageReactionRemoveEventHandler)
	s.AddHandler(MessageCreateEventHandler)
	s.AddHandler(ChannelDeleteEventHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuildEmojis

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

func Stop() {
	log.Println("Closing bot connection")
	s.Close()
	log.Println("Closing database connection")
	a.DB.Close()
}
