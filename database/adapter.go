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
	return &a
}

func (a *Adapter) InitServer() {

}

func (a *Adapter) CutServer() {

}

func (a *Adapter) ResetServer() {

}

func (a *Adapter) GetGuildWelcome() {

}

func (a *Adapter) SetGuildWelcomeGuildId() {

}

func (a *Adapter) SetGuildWelcomeImage() {

}

func (a *Adapter) SetGuildWelcomeImageText() {

}

func (a *Adapter) SetGuildWelcomeMessageText() {

}
