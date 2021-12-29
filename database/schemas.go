package database

type GuildWelcome struct {
	GuildID     string
	ChannelID   string
	Type        string
	MessageText string
	Image       string
	ImageText   string
}

func NewDefaultGuildWelcome() *GuildWelcome {
	return &GuildWelcome{
		ChannelID:   "",
		Type:        "image",
		MessageText: "hi %mention% <3 welcome to %guild% :)",
		Image:       "original",
		ImageText:   "%username% joined the server",
	}
}
