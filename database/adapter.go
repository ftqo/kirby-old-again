package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Adapter struct {
	pool *pgxpool.Pool
}

func Open() *Adapter {
	dbhost := os.Getenv("DBHOST")
	dbport, err := strconv.Atoi(os.Getenv("DBPORT"))
	if err != nil {
		log.Panicf("failed to parse database port environment variable: %v", err)
	}
	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	database := os.Getenv("DATABASE")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, database)

	p, err := pgxpool.Connect(context.TODO(), dsn)
	if err != nil {
		log.Panicf("failed to open connextion pool: %v", err)
	}
	err = p.Ping(context.TODO())
	if err != nil {
		log.Panicf("failed to connect to connection pool: %v", err)
	}
	log.Print("connected to database !")
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
		log.Panicf("failed to acquire connection for database initialization: %v", err)
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
		log.Printf("failed to execute statement for initialize database: %v", err)
	}
}

func (a *Adapter) InitGuild(guildId string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for guild initialization: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for guild initialization: %v", err)
	}

	dgw := NewDefaultGuildWelcome()

	statement := `
	INSERT INTO guild_welcome (guild_id, channel_id, type, message_text, image, image_text)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (guild_id) DO NOTHING`
	_, err = tx.Exec(context.TODO(), statement, guildId, dgw.ChannelID, dgw.Type, dgw.Text, dgw.Image, dgw.ImageText)
	if err != nil {
		log.Printf("failed to execute statement for guild initialization: %v", err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("failed to commit transaction for guild initialization: %v", err)
	}
}

func (a *Adapter) CutGuild(guildId string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for guild cutting: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for guild cutting: %v", err)
	}

	statement := `
	DELETE FROM guild_welcome WHERE guild_id = $1`
	_, err = tx.Exec(context.TODO(), statement, guildId)
	if err != nil {
		log.Printf("failed to execute statement for guild cutting: %v", err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("failed to commit transaction for guild cutting: %v", err)
	}
}

func (a *Adapter) ResetServer(guildId string) {
	a.CutGuild(guildId)
	a.InitGuild(guildId)
}

func (a *Adapter) GetGuildWelcome(guildId string) GuildWelcome {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for guild welcome: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for guild welcome: %v", err)
	}

	statement := `
	SELECT guild_id, channel_id, type, message_text, image, image_text FROM guild_welcome WHERE guild_id = $1`
	row := tx.QueryRow(context.TODO(), statement, guildId)
	gw := GuildWelcome{}
	err = row.Scan(&gw.GuildID, &gw.ChannelID, &gw.Type, &gw.Text, &gw.Image, &gw.ImageText)
	if err != nil {
		log.Printf("failed to scan query for guild welcome: %v", err)
	}

	return gw
}

func (a *Adapter) SetGuildWelcomeChannel(guildId, channelId string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for set guild welcome channel: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for set guild welcome channel: %v", err)
	}

	statement := `
	UPDATE guild_welcome SET channel_id = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, channelId, guildId)
	if err != nil {
		log.Printf("failed to execute statement for set guild welcome channel: %v", err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("failed to commit transaction for set guild welcome channel: %v", err)
	}
}

func (a *Adapter) SetGuildWelcomeText(guildId, messageText string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for set guild welcome message text: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for set guild welcome message text: %v", err)
	}

	statement := `
	UPDATE guild_welcome SET message_text = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, messageText, guildId)
	if err != nil {
		log.Printf("failed to execute statement for set guild welcome message text: %v", err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("failed to commit transaction for set guild welcome message text: %v", err)
	}
}

func (a *Adapter) SetGuildWelcomeImage(guildId, image string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for set guild welcome image: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for set guild welcome image: %v", err)
	}

	statement := `
	UPDATE guild_welcome SET image = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, image, guildId)
	if err != nil {
		log.Printf("failed to execute statement for set guild welcome image: %v", err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("failed to commit transaction for set guild welcome image: %v", err)
	}
}

func (a *Adapter) SetGuildWelcomeImageText(guildId, imageText string) {
	conn, err := a.pool.Acquire(context.TODO())
	if err != nil {
		log.Printf("failed to acquire connection for set guild welcome image text: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Printf("failed to begin transaction for set guild welcome image text: %v", err)
	}

	statement := `
	UPDATE guild_welcome SET image_text = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, imageText, guildId)
	if err != nil {
		log.Printf("failed to execute statement for set guild welcome image text: %v", err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("failed to commit transaction for set guild welcome image text: %v", err)
	}
}
