package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Adapter struct {
	pool *pgxpool.Pool
}

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "jesus"
	database = "testing"
)

func Open() *Adapter {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	p, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Panicf("Could not open connextion pool: %v", err)
	}
	err = p.Ping(context.Background())
	if err != nil {
		log.Panicf("Could not connect to connection pool: %v", err)
	}
	log.Println("Database connected!")
	a := Adapter{
		pool: p,
	}
	a.initDatabase()
	return &a
}

func (a *Adapter) Close() {
	a.pool.Close()
}

func (a *Adapter) initDatabase() {
	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Panicf("Could not acquire connection for database initialization: %v", err)
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
	_, err = conn.Exec(context.Background(), statement)
	if err != nil {
		log.Printf("Could not execute statement for initialize database: %v", err)
	}
}

func (a *Adapter) InitGuild(guildId string) {
	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Could not acquire connection for guild initialization: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Could not begin transaction for guild initialization: %v", err)
	}

	dgw := NewDefaultGuildWelcome()

	statement := `
	INSERT INTO guild_welcome (guild_id, channel_id, type, message_text, image, image_text)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (guild_id) DO NOTHING`
	_, err = tx.Exec(context.Background(), statement, guildId, dgw.ChannelID, dgw.Type, dgw.MessageText, dgw.Image, dgw.ImageText)
	if err != nil {
		log.Printf("Could not execute statement for guild initialization: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Could not commit transaction for guild initialization: %v", err)
	}
}

func (a *Adapter) CutGuild(guildId string) {
	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Could not acquire connection for guild cutting: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Could not begin transaction for guild cutting: %v", err)
	}

	statement := `
	DELETE FROM guild_welcome WHERE guild_id = $1`
	_, err = tx.Exec(context.Background(), statement, guildId)
	if err != nil {
		log.Printf("Could not execute statement for guild cutting: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Could not commit transaction for guild cutting: %v", err)
	}
}

func (a *Adapter) ResetServer(guildId string) {
	a.CutGuild(guildId)
	a.InitGuild(guildId)
}

func (a *Adapter) GetGuildWelcome(guildId string) GuildWelcome {
	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Could not acquire connection for guild welcome: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Could not begin transaction for guild welcome: %v", err)
	}

	statement := `
	SELECT guild_id, channel_id, type, message_text, image, image_text FROM guild_welcome WHERE guild_id = $1`
	row := tx.QueryRow(context.Background(), statement, guildId)
	gw := GuildWelcome{}
	err = row.Scan(&gw.GuildID, &gw.ChannelID, &gw.Type, &gw.MessageText, &gw.Image, &gw.ImageText)
	if err != nil {
		log.Printf("Could not scan query for guild welcome: %v", err)
	}

	return gw
}

func (a *Adapter) SetGuildWelcomeChannel(guildId, channelId string) {

}

func (a *Adapter) SetGuildWelcomeMessageText(guildId, messageText string) {

}

func (a *Adapter) SetGuildWelcomeImage(guildId, image string) {

}

func (a *Adapter) SetGuildWelcomeImageText(guildId, imageText string) {

}
