package discord

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/logger"
)

const (
	PfpSize = 256

	width   = 848
	height  = 477
	margin  = 15
	res     = 1
	titlSz  = 100
	stitlSz = 80
)

type welcomeMessageInfo struct {
	mention   string
	nickname  string
	username  string
	guildName string
	avatarURL string
	members   int
}

func generateWelcomeMessage(gw database.GuildWelcome, wi welcomeMessageInfo) discordgo.MessageSend {
	var msg discordgo.MessageSend

	r := strings.NewReplacer("%mention%", wi.mention, "%nickname%", wi.nickname, "%username%", wi.username, "%guild%", wi.guildName)
	gw.Text = r.Replace(gw.Text)
	gw.ImageText = r.Replace(gw.ImageText)

	msg.Content = gw.Text

	switch gw.Type {
	case "embed":
		logger.L.Error().Msg("Embedded welcome messages not implemented, sending plain")
	case "image":
		ctx := gg.NewContextForImage(assets.Images[gw.Image])
		resp, err := http.Get(wi.avatarURL)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to get avatar URL")
		}
		defer resp.Body.Close()

		pfpBuf := bytes.Buffer{}
		_, err = io.Copy(&pfpBuf, resp.Body)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to copy pfp to bytes buffer")
		}
		rawPfp, _, err := image.Decode(&pfpBuf)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to decode profile picture")
		}
		var pfp image.Image
		if rawPfp.Bounds().Max.X != PfpSize {
			pfp = image.Image(transform.Resize(rawPfp, PfpSize, PfpSize, transform.Linear))
		} else {
			pfp = rawPfp
		}

		ctx.SetColor(color.RGBA{52, 45, 50, 130})
		ctx.DrawRectangle(margin, margin, width-(2*margin), height-(2*margin))
		ctx.Fill()
		ctx.ClearPath()

		ctx.SetColor(color.White)
		ctx.DrawCircle(width/2, height*(44.0/100.0), PfpSize/2+3)
		ctx.SetLineWidth(5)
		ctx.Stroke()
		ctx.DrawCircle(width/2, height*(44.0/100.0), PfpSize/2)
		ctx.Clip()
		ctx.DrawImage(pfp, width/2-PfpSize/2, height*44/100-PfpSize/2)
		ctx.ResetClip()

		fontLarge := assets.Fonts["coolveticaLarge"]
		fontSmall := assets.Fonts["coolveticaSmall"]

		ctx.SetFontFace(fontLarge)
		ctx.DrawStringAnchored(gw.ImageText, width/2, height*78/100, 0.5, 0.5)
		ctx.SetFontFace(fontSmall)
		ctx.DrawStringAnchored("member #"+strconv.Itoa(wi.members), width/2, height*85/100, 0.5, 0.5)

		buf := bytes.Buffer{}
		err = png.Encode(&buf, ctx.Image())
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to encode image into bytes buffer")
		}

		f := &discordgo.File{
			Name:        "welcome_" + wi.nickname + ".jpg",
			ContentType: "image/jpeg",
			Reader:      &buf,
		}
		msg.Files = append(msg.Files, f)
	}
	return msg
}
