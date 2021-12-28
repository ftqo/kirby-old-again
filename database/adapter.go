package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Adapter struct {
	DB *pgxpool.Pool
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
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Panicln(err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Database connected!")
	a := Adapter{
		DB: conn,
	}
	a.prepareDatabase()
	return &a
}

func (a *Adapter) InitServer(guildId string) {

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
