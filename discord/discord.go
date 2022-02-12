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
var tg string
var cc []*discordgo.ApplicationCommand

func Start(token, testGuild string) {
	tg = testGuild
	var err error
	h = new(files.Hoarder)
	h.LoadFiles()
	a = database.Open()

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("failed to initialize discordgo session: %v", err)
	}
	s.AddHandler(ReadyHandler)
	s.AddHandler(GuildCreateEventHandler)
	s.AddHandler(GuildDeleteEventHandler)
	s.AddHandler(GuildMemberAddEventHandler)
	s.AddHandler(ChannelDeleteEventHandler)
	s.AddHandler(InteractionCreateEventHandler)

	s.Identify.Intents = discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages

	err = s.Open()
	if err != nil {
		log.Fatalf("failed to open the discord session: %v", err)
	}
}

func Stop() {
	log.Print("removing bot commands !")
	for _, cmd := range cc {
		err := s.ApplicationCommandDelete(s.State.User.ID, tg, cmd.ID)
		if err != nil {
			log.Printf("failed to delete %q command: %v", cmd.Name, err)
		}
	}
	log.Print("closing bot connection !")
	s.Close()
	log.Print("closing database connection !")
	a.Close()
}
