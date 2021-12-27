package welcome

import (
	"bytes"
	"image/jpeg"
	"log"

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

func GenerateWelcome(w database.GuildWelcome, u *discordgo.User) (discordgo.MessageSend, string) {
	var msg discordgo.MessageSend
	switch w.Type {
	case "plain":
		msg.Content = w.MessageText
	case "embed":
		log.Println("not implemented")
	case "image":
		background := softwarebackend.New(width, height)
		ctx := canvas.New(background)

		_, err := ctx.LoadFont("")
		if err != nil {
			log.Panic(err)
		}

		img := ctx.GetImageData(0, 0, width, height)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, nil)
		if err != nil {
			log.Fatalln(err)
		}
		name := "welcome_" + u.Username + ".jpg"
		cont := "image/jpeg"
		read := bytes.NewReader(buf.Bytes())
		f := discordgo.File{
			Name:        name,
			ContentType: cont,
			Reader:      read,
		}
		msg.Files = append(msg.Files, &f)
	}
	return msg, w.Channel
}
