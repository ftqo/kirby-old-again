package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ftqo/kirby/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Adapter struct {
	pool *pgxpool.Pool
}

func Open() *Adapter {
	dbhost := os.Getenv("DB_HOST")
	dbport, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.L.Panic().Err(err).Msg("failed to parse database port environment variable")
	}
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASS")
	database := os.Getenv("DB_DATABASE")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, database)

	p, err := pgxpool.Connect(context.TODO(), dsn)
	if err != nil {
		logger.L.Panic().Err(err).Msg("failed to open connextion pool")
	}
	err = p.Ping(context.TODO())
	if err != nil {
		logger.L.Panic().Err(err).Msg("failed to connect to connection pool")
	}
	logger.L.Info().Msg("connected to database !")
	a := &Adapter{
		pool: p,
	}
	a.initDatabase()
	return a
}

func (a *Adapter) Close() {
	a.pool.Close()
}

func (a *Adapter) initDatabase() {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Panic().Err(err).Msg("failed to acquire connection for database initialization")
	}
	defer conn.Release()
	statement := `
	CREATE TABLE IF NOT EXISTS guild_welcome (
		guild_id TEXT PRIMARY KEY,
		channel_id TEXT NOT NULL,
		type TEXT NOT NULL,
		message_text TEXT NOT NULL,
		image TEXT NOT NULL,
		image_text TEXT NOT NULL
	);`
	_, err = conn.Exec(context.TODO(), statement)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for initialize database")
	}
}

func (a *Adapter) InitGuild(guildId string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for guild initialization")
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for guild initialization")
	}
	dgw := NewDefaultGuildWelcome()
	statement := `
	INSERT INTO guild_welcome (guild_id, channel_id, type, message_text, image, image_text)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (guild_id) DO NOTHING`
	_, err = tx.Exec(context.TODO(), statement, guildId, dgw.ChannelID, dgw.Type, dgw.Text, dgw.Image, dgw.ImageText)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for guild initialization")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for guild initialization")
	}
}

func (a *Adapter) CutGuild(guildId string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for guild cutting")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for guild cutting")
	}
	statement := `
	DELETE FROM guild_welcome WHERE guild_id = $1`
	_, err = tx.Exec(context.TODO(), statement, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for guild cutting")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for guild cutting")
	}
}

func (a *Adapter) ResetGuild(guildId string) {
	a.CutGuild(guildId)
	a.InitGuild(guildId)
}

func (a *Adapter) GetGuildWelcome(guildId string) GuildWelcome {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for guild welcome")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for guild welcome")
	}
	statement := `
	SELECT guild_id, channel_id, type, message_text, image, image_text FROM guild_welcome WHERE guild_id = $1`
	row := tx.QueryRow(context.TODO(), statement, guildId)
	gw := GuildWelcome{}
	err = row.Scan(&gw.GuildID, &gw.ChannelID, &gw.Type, &gw.Text, &gw.Image, &gw.ImageText)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to scan query for guild welcome")
	}
	return gw
}

func (a *Adapter) SetGuildWelcomeChannel(guildId, channelId string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for set guild welcome channel")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for set guild welcome channel")
	}
	statement := `
	UPDATE guild_welcome SET channel_id = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, channelId, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for set guild welcome channel")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for set guild welcome channel")
	}
}

func (a *Adapter) SetGuildWelcomeType(guildId, welcomeType string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for set guild welcome image")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for set guild welcome image")
	}
	statement := `
	UPDATE guild_welcome SET type = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, welcomeType, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for set guild welcome image")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for set guild welcome image")
	}
}

func (a *Adapter) SetGuildWelcomeText(guildId, messageText string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for set guild welcome message text")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for set guild welcome message text")
	}
	statement := `
	UPDATE guild_welcome SET message_text = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, messageText, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for set guild welcome message text")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for set guild welcome message text")
	}
}

func (a *Adapter) SetGuildWelcomeImage(guildId, image string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for set guild welcome image")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for set guild welcome image")
	}
	statement := `
	UPDATE guild_welcome SET image = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, image, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for set guild welcome image")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for set guild welcome image")
	}
}

func (a *Adapter) SetGuildWelcomeImageText(guildId, imageText string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to acquire connection for set guild welcome image text")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to begin transaction for set guild welcome image text")
	}
	statement := `
	UPDATE guild_welcome SET image_text = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, imageText, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to execute statement for set guild welcome image text")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("failed to commit transaction for set guild welcome image text")
	}
}
