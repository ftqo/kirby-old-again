package discord

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/files"
)

var s *discordgo.Session
var assets *files.Assets
var adapter *database.Adapter
var cc []*discordgo.ApplicationCommand
var tg string
var rmCmd bool

func Start(token, testGuild string, rmCommands string) {
	var err error
	tg = testGuild
	rmCmd, err = strconv.ParseBool(rmCommands)
	if err != nil {
		rmCmd = false
	}
	assets = files.GetAssets()
	adapter = database.Open()

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
	cc, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, tg, commands)
	if err != nil {
		log.Panicf("failed to create command application commands: %v", err)
	}
	log.Print("loaded slash commands !")
}

func Stop() {
	if rmCmd {
		log.Print("removing bot commands !")
		for _, c := range cc {
			err := s.ApplicationCommandDelete(s.State.User.ID, tg, c.ID)
			if err != nil {
				log.Printf("failed to delete %q command: %v", c.Name, err)
			}
		}
	}
	log.Print("closing bot connection !")
	s.Close()
	log.Print("closing database connection !")
	adapter.Close()
}
