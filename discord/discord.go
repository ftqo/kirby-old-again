package discord

import (
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/logger"
)

var s *discordgo.Session

var (
	cc     []*discordgo.ApplicationCommand
	tg     string
	rmCmds bool
)

func Start(token, testGuild string, rmCommands string) {
	var err error
	tg = testGuild
	rmCmds, err = strconv.ParseBool(rmCommands)
	if err != nil {
		rmCmds = false
	}
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to initialize discordgo session")
	}

	s.AddHandler(readyHandler)
	s.AddHandler(guildCreateEventHandler)
	s.AddHandler(guildDeleteEventHandler)
	s.AddHandler(guildMemberAddEventHandler)
	s.AddHandler(channelDeleteEventHandler)
	s.AddHandler(interactionCreateEventHandler)
	s.Identify.Intents = discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages
	err = s.Open()
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to open the Discord session")
	}
	logger.L.Info().Msg("Opened connection with Discord")
	cc, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, tg, commands)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to create application commands")
	}
	logger.L.Info().Msg("Updated applications commands")
}

func Stop() {
	if rmCmds {
		for _, c := range cc {
			err := s.ApplicationCommandDelete(s.State.User.ID, tg, c.ID)
			if err != nil {
				logger.L.Error().Err(err).Msgf("Failed to delete application command %s", c.Name)
			}
		}
		logger.L.Info().Msg("Deleted application commands as set by RMCMDS")
	}
	err := s.CloseWithCode(websocket.CloseNormalClosure)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to close bot connection properly")
	} else {
		logger.L.Info().Msg("Closed connection with Discord")
	}
}
