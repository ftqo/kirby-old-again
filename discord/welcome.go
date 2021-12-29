package discord

import (
	"bytes"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fittsqo/kirby/database"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	width  = 1696
	height = 954
	pfp    = 512
	margin = 30
)

type welcomeMessageInfo struct {
	mention   string
	nickname  string
	username  string
	guildname string
}

func GenerateWelcomeMessage(gw database.GuildWelcome, wi welcomeMessageInfo) discordgo.MessageSend {
	var msg discordgo.MessageSend

	r := strings.NewReplacer("%mention%", wi.mention, "%nickname", wi.nickname, "%username%", wi.username, "%guild%", wi.guildname)
	gw.MessageText = r.Replace(gw.MessageText)
	gw.ImageText = r.Replace(gw.ImageText)

	msg.Content = gw.MessageText

	switch gw.Type {
	case "plain":
		log.Println("Generating plain welcome message")
		msg.Content = gw.MessageText
	case "embed":
		log.Println("Embedded welcome messages not implemented, sending plain")
		msg.Content = gw.MessageText
	case "image":
		log.Println("Generating image welcome message")

		cv := canvas.New(width, height)
		ctx := canvas.NewContext(cv)
		ctx.DrawImage(0, 0, h.Images[gw.Image], 1.0)
		buf := &bytes.Buffer{}
		cw := renderers.PNG()
		cw(buf, cv)

		f := discordgo.File{
			Name:        "welcome_" + wi.nickname + ".jpg",
			ContentType: "image/jpeg",
			Reader:      bytes.NewReader(buf.Bytes()),
		}

		msg.Files = append(msg.Files, &f)
	}
	return msg
}
