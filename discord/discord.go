package discord

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/files"
	"github.com/ftqo/kirby/logger"

	"github.com/gorilla/websocket"
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
		logger.L.Panic().Err(err).Msg("Failed to initialize discordgo session")
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
		logger.L.Panic().Err(err).Msg("Failed to open the discord session")
	}
	cc, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, tg, commands)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to create command application commands")
	}
	logger.L.Info().Msg("Loaded slash commands")
}

func Stop() {
	if rmCmd {
		logger.L.Info().Msg("Removing bot commands as enabled")
		for _, c := range cc {
			err := s.ApplicationCommandDelete(s.State.User.ID, tg, c.ID)
			if err != nil {
				logger.L.Error().Err(err).Msgf("Failed to delete %q command", c.Name)
			}
		}
	}
	logger.L.Info().Msg("Closing bot connection")
	err := s.CloseWithCode(websocket.CloseNormalClosure)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to close with code restart")
	}
	logger.L.Info().Msg("Closing database connection")
	adapter.Close()
}
