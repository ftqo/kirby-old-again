package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/logger"
)

func readyHandler(s *discordgo.Session, e *discordgo.Ready) {
	logger.L.Debug().Msg("[READY]")
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
		logger.L.Panic().Err(err).Msg("Failed to update status")
	}
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, tg, commands)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to create application commands")
	}
	logger.L.Info().Msg("Updated applications commands")
}

func resumeHandler(s *discordgo.Session, e *discordgo.Resumed) {
	logger.L.Debug().Msg("[RESUMED]")
}

func guildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot turns on or joins a guild
	logger.L.Debug().Msgf("[GUILD_CREATE] %s (%s)", e.Guild.Name, e.Guild.ID)
	database.InitGuild(e.Guild.ID)
}

func guildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild
	logger.L.Debug().Msgf("[GUILD_DELETE] %s (%s) (KICKED: %t)", e.Guild.Name, e.Guild.ID, !e.Unavailable)
	if !e.Unavailable {
		database.CutGuild(e.Guild.ID)
	}
}

func guildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	logger.L.Debug().Msgf("[GUILD_MEMBER_ADD] %s (%s) JOINED %s", e.User.String(), e.User.ID, e.GuildID)
	g, err := s.State.Guild(e.GuildID)
	if err != nil {
		logger.L.Warn().Err(err).Msgf("Failed to get guild from cache when GuildMemberAdd was fired")
		g, err = s.Guild(e.GuildID)
		if err != nil {
			logger.L.Error().Err(err).Msgf("Failed to get guild from direct request")
			return
		}
		s.State.GuildAdd(g)
	}
	gw := database.GetGuildWelcome(g.ID)
	if gw.ChannelID != "" {
		wi := welcomeMessageInfo{
			mention:   e.User.Mention(),
			nickname:  e.User.Username,
			username:  e.User.String(),
			guildName: g.Name,
			avatarURL: e.User.AvatarURL(fmt.Sprint(PfpSize)),
			members:   g.MemberCount,
		}
		welcome := generateWelcomeMessage(gw, wi)
		_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
		if err != nil {
			logger.L.Error().Err(err).Msgf("Failed to send welcome message")
		}
	}
}

func channelDeleteEventHandler(s *discordgo.Session, e *discordgo.ChannelDelete) {
	logger.L.Debug().Msgf("[CHANNEL_DELETE] %s IN %s", e.Channel.Name, e.GuildID)
	gw := database.GetGuildWelcome(e.GuildID)
	if e.Channel.ID == gw.ChannelID {
		database.SetGuildWelcomeChannel(e.GuildID, "")
	}
}

func interactionCreateEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger.L.Debug().Msgf("[INTERACTION_CREATE] %s (%s) USED %s", i.Member.User.String(), i.Member.User.ID, i.ApplicationCommandData().Name)
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
			h(s, i)
		}
	}
}
