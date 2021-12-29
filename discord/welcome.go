package discord

import (
	"bytes"
	"image/jpeg"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fittsqo/kirby/database"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/softwarebackend"
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
		background := softwarebackend.New(width, height)
		ctx := canvas.New(background)

		i, err := ctx.LoadImage(h.Images[0])
		if err != nil {
			log.Panicln(err)
		}
		ctx.DrawImage(i)

		_, err = ctx.LoadFont(h.Font)
		if err != nil {
			log.Panic(err)
		}

		img := ctx.GetImageData(0, 0, width, height)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, nil)
		if err != nil {
			log.Fatalln(err)
		}
		f := discordgo.File{
			Name:        "welcome_" + wi.nickname + ".jpg",
			ContentType: "image/jpeg",
			Reader:      bytes.NewReader(buf.Bytes()),
		}
		msg.Files = append(msg.Files, &f)
	}
	return msg
}
