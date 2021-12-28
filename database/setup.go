package database

import (
	"context"
	"log"
)

const createWelcomeTable = `CREATE TABLE IF NOT EXISTS guild_welcome (
	guild_id TEXT PRIMARY KEY,
	type SMALLINT NOT NULL,
	channel TEXT NOT NULL,
	image SMALLINT NOT NULL,
	image_text TEXT NOT NULL,
	message_text TEXT NOT NULL
);`

func (a *Adapter) prepareDatabase() {
	tx, err := a.DB.Begin(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	_, err = tx.Exec(context.Background(), createWelcomeTable)
	if err != nil {
		log.Panicln(err)
	}
	tx.Commit(context.Background())
}
