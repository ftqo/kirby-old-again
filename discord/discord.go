package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/files"
)

var s *discordgo.Session
var h *files.Hoarder
var a *database.Adapter

func Start(token string) {
	var err error
	h = new(files.Hoarder)
	h.LoadFiles()
	a = database.Open()

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}
	s.AddHandler(ReadyHandler)
	s.AddHandler(GuildCreateEventHandler)
	s.AddHandler(GuildDeleteEventHandler)
	s.AddHandler(GuildMemberAddEventHandler)
	s.AddHandler(ChannelDeleteEventHandler)

	s.AddHandler(MessageCreateEventHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildEmojis

	err = s.Open()
	if err != nil {
		log.Fatalf("failed to open the discord session: %v", err)
	}
}

func Stop() {
	log.Println("closing bot connection !")
	s.Close()
	log.Println("closing database connection !")
	a.Close()
}
