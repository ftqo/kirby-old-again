package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() {
	flag.Parse()
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(readyHandler)
	s.AddHandler(guildCreateEventHandler)
	s.AddHandler(guildDeleteEventHandler)
	s.AddHandler(guildMemberAddEventHandler)
	s.AddHandler(guildMemberRemoveEventHandler)
	s.AddHandler(messageDeleteEventHandler)
	s.AddHandler(messageReactionAddEventHandler)
	s.AddHandler(messageReactionRemoveEventHandler)
	s.AddHandler(messageCreateEventHandler)
	s.AddHandler(channelDeleteEventHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers

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
