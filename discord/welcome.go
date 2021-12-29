package discord

import (
	"bytes"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/bwmarrin/discordgo"
	"github.com/fittsqo/kirby/database"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	width   = 848
	height  = 477
	pfpSize = 256
	margin  = 15
	res     = 1
)

type welcomeMessageInfo struct {
	mention   string
	nickname  string
	username  string
	guildName string
	avatarURL string
}

func GenerateWelcomeMessage(gw database.GuildWelcome, wi welcomeMessageInfo) discordgo.MessageSend {
	var msg discordgo.MessageSend

	r := strings.NewReplacer("%mention%", wi.mention, "%nickname", wi.nickname, "%username%", wi.username, "%guild%", wi.guildName)
	gw.MessageText = r.Replace(gw.MessageText)
	gw.ImageText = r.Replace(gw.ImageText)

	msg.Content = gw.MessageText

	switch gw.Type {
	case "plain":
		log.Println("Generating plain welcome message")
	case "embed":
		log.Println("Embedded welcome messages not implemented, sending plain")
	case "image":
		log.Println("Generating image welcome message")
		cv := canvas.New(width, height)
		ctx := canvas.NewContext(cv)
		resp, err := http.Get(wi.avatarURL)
		if err != nil {
			log.Printf("Could not get avatar URL: %v", err)
		}
		defer resp.Body.Close()

		pfpBuf := &bytes.Buffer{}
		_, err = io.Copy(pfpBuf, resp.Body)
		if err != nil {
			log.Printf("Could not copy pfp to buffer: %v", err)
		}
		rawPfp, _, err := image.Decode(pfpBuf)
		if err != nil {
			log.Printf("Could not decode profile picture: %v", err)
		}
		var pfp image.Image
		if rawPfp.Bounds().Max.X != pfpSize {
			pfp = image.Image(transform.Resize(rawPfp, pfpSize, pfpSize, transform.Linear))
		} else {
			pfp = rawPfp
		}

		ctx.DrawImage(0, 0, h.Images[gw.Image], res)
		ctx.DrawImage(0, 0, pfp, res)

		buf := &bytes.Buffer{}
		cw := renderers.JPEG()
		cw(buf, cv)
		f := &discordgo.File{
			Name:        "welcome_" + wi.nickname + ".jpg",
			ContentType: "image/jpeg",
			Reader:      buf,
		}
		msg.Files = append(msg.Files, f)
	}
	return msg
}
