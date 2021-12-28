package database

type GuildWelcome struct {
	GuildId     string `db:"guild_id"`
	Type        string `db:"type"`
	Channel     string `db:"channel"`
	Image       int    `db:"image"`
	ImageText   string `db:"image_text"`
	MessageText string `db:"message_text"`
}

// type ReactionMessag struct {
// }
