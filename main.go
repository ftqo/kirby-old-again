package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/fittsqo/kirby/services/discord"
	"github.com/fittsqo/kirby/services/files"
)

var (
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session
var f *files.HoardedFiles

func init() {
	flag.Parse()
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	f = new(files.HoardedFiles)
	f.LoadImages()

	s.AddHandler(discord.ReadyHandler)
	s.AddHandler(discord.GuildCreateEventHandler)
	s.AddHandler(discord.GuildDeleteEventHandler)
	s.AddHandler(discord.GuildMemberAddEventHandler)
	s.AddHandler(discord.GuildMemberRemoveEventHandler)
	s.AddHandler(discord.MessageDeleteEventHandler)
	s.AddHandler(discord.MessageReactionAddEventHandler)
	s.AddHandler(discord.MessageReactionRemoveEventHandler)
	s.AddHandler(discord.MessageCreateEventHandler)
	s.AddHandler(discord.ChannelDeleteEventHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuildEmojis

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}
