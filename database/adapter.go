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

func Connect() *Adapter {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	p, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Panicln(err)
	}
	err = p.Ping(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Database connected!")
	a := Adapter{
		pool: p,
	}
	a.prepareDatabase()
	return &a
}

func (a *Adapter) Close() {
	a.pool.Close()
}

func (a *Adapter) prepareDatabase() {
	statement := `
	CREATE TABLE IF NOT EXISTS guild_welcome (
		guild_id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		channel TEXT NOT NULL,
		message_text TEXT NOT NULL,
		image TEXT NOT NULL,
		image_text TEXT NOT NULL
	);`

	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), statement)
	if err != nil {
		log.Panicln(err)
	}
}

func (a *Adapter) InitServer(guildId string) {
	statement := `
	INSERT INTO guild_welcome (guild_id, type, channel, message_text, image, image_text)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (guild_id) DO NOTHING`

	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	tx.Prepare(context.Background(), "init", statement)
	dgw := NewDefaultGuildWelcome()
	_, err = tx.Exec(context.Background(), "init", guildId, dgw.Type, dgw.Channel, dgw.MessageText, dgw.Image, dgw.ImageText)
	if err != nil {
		log.Fatalln(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}

func (a *Adapter) CutServer(guildId string) {

}

func (a *Adapter) ResetServer(guildId string) {
	a.CutServer(guildId)
	a.InitServer(guildId)
}

func (a *Adapter) GetGuildWelcome(guildId string) *GuildWelcome {
	gw := GuildWelcome{}
	return &gw
}

func (a *Adapter) SetGuildWelcomeChannel(guildId, channelId string) {

}

func (a *Adapter) SetGuildWelcomeMessageText(guildId, messageText string) {

}

func (a *Adapter) SetGuildWelcomeImage(guildId, image string) {

}

func (a *Adapter) SetGuildWelcomeImageText(guildId, imageText string) {

}
