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
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == "879505042200227840" || m.Author.ID == "875995595465171004" { // bot IDs
			return
		}
		if m.Content == "!ping" {
			_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	s.AddHandler(func(s *discordgo.Session, a *discordgo.GuildMemberAdd) {
		channelID := "922379294477529108"
		message := "hi"
		_, err := s.ChannelMessageSend(channelID, message)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Attempted to send %s to channel %s\n", message, channelID)
	})

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
