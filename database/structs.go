package database

type GuildWelcome struct {
	GuildId     string `db:"guild_id"`
	Channel     string `db:"channel"`
	Type        string `db:"type"`
	Image       int    `db:"image"`
	ImageText   string `db:"image_text"`
	MessageText string `db:"message_text"`
}

// type ReactionMessag struct {
// }
