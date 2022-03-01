package discord

import (
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/logger"
)

var (
	s  *discordgo.Session
	tg string
)

func Start(token, sessionID, sequence, testGuild string) {
	var err error // prevent shadowing with :=
	tg = testGuild
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to initialize local dgo session")
	}
	s.AddHandler(readyHandler)
	s.AddHandler(resumeHandler)
	s.AddHandler(guildCreateEventHandler)
	s.AddHandler(guildDeleteEventHandler)
	s.AddHandler(guildMemberAddEventHandler)
	s.AddHandler(channelDeleteEventHandler)
	s.AddHandler(interactionCreateEventHandler)
	s.Identify.Intents = discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages
	if sessionID != "" && sequence != "" { // if sessionID and sequence are set, will attempt to resume connection
		s.SessionID = sessionID
		seq, err := strconv.Atoi(sequence)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to convert sequence to int")
		}
		seq64 := int64(seq)
		s.Sequence = &seq64
	}
	err = s.Open()
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to open the Discord session")
	}
	logger.L.Info().Msg("Opened connection with Discord")
}

func Stop() (string, string) {
	err := s.CloseWithCode(websocket.CloseServiceRestart)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to close bot connection properly")
	} else {
		logger.L.Info().Msg("Closed connection with Discord")
	}
	return s.SessionID, strconv.Itoa(int(*s.Sequence))
}
