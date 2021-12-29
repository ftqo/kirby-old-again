package database

type GuildWelcome struct {
	GuildId     string
	Type        string
	Channel     string
	MessageText string
	Image       string
	ImageText   string
}

func NewDefaultGuildWelcome() *GuildWelcome {
	dgw := GuildWelcome{
		Type:        "image",
		Channel:     "",
		MessageText: "hi %user_mention% <3 welcome to %%guild_name% :)",
		Image:       "original",
		ImageText:   "%user_tag% joined the server",
	}
	return &dgw
}
