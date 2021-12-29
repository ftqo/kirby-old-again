package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	usd := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "the stars",
				Type: 3,
			},
		},
	}
	err := s.UpdateStatusComplex(usd)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Bot connected!")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot turns on or joins a guild
	a.InitServer(e.Guild.ID)
}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild
	if !e.Unavailable {
		a.CutServer(e.Guild.ID)
	}
}

func GuildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {

}

func MessageCreateEventHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Content == "!simwelcome" {
		gw := a.GetGuildWelcome(e.GuildID) // grab info from database
		g, err := s.Guild(e.GuildID)
		if err != nil {
			log.Fatalln(err)
		}
		wi := welcomeMessageInfo{ // generate replacements for placeholders
			"<@" + e.Author.ID + ">",
			e.Author.Username,
			e.Author.Username + "#" + e.Author.Discriminator,
			g.Name,
		}
		msg := GenerateWelcomeMessage(gw, wi)
		_, err = s.ChannelMessageSendComplex(e.ChannelID, &msg) // change e.ChannelID to gw.Channel
		if err != nil {
			log.Fatalln(err)
		}
	}
}
