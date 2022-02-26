package files

import (
	"bytes"
	"image"
	"io/ioutil"
	"path"
	"runtime"
	"strings"

	"github.com/ftqo/kirby/logger"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Assets struct {
	Fonts  map[string]font.Face
	Images map[string]image.Image
}

func GetAssets() *Assets {
	var ass Assets
	ass.Images = make(map[string]image.Image)
	ass.Fonts = make(map[string]font.Face)
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	assetsPath := path.Join(d, "../assets")
	imagesPath := path.Join(assetsPath, "images")
	fontsPath := path.Join(assetsPath, "fonts")
	imgs, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		logger.L.Panic().Err(err).Msgf("Failed to read directory %s", imagesPath)
	}
	for _, file := range imgs {
		fp := path.Join(imagesPath, file.Name())
		bts, err := ioutil.ReadFile(fp)
		if err != nil {
			logger.L.Panic().Err(err).Msgf("Failed to read file %s", fp)
		}
		fn := file.Name()
		noPre := fn[strings.LastIndex(file.Name(), "-")+1:]
		noExt := noPre[:strings.Index(noPre, ".")]
		ass.Images[noExt], _, err = image.Decode(bytes.NewReader(bts))
		if err != nil {
			logger.L.Panic().Err(err).Msgf("Failed to decode %s", fp)
		}
		logger.L.Info().Msgf("Loaded %s", fp)
	}
	fonts, err := ioutil.ReadDir(fontsPath)
	if err != nil {
		logger.L.Panic().Err(err).Msgf("Failed to read directory %s", fonts)
	}
	for _, file := range fonts {
		fp := path.Join(fontsPath, file.Name())
		bts, err := ioutil.ReadFile(fp)
		if err != nil {
			logger.L.Panic().Err(err).Msgf("Failed to read file %s", fp)
		}
		fn := file.Name()
		noPre := fn[strings.LastIndex(file.Name(), "-")+1:]
		noExt := noPre[:strings.Index(noPre, ".")]
		fnt, err := truetype.Parse(bts)
		if err != nil {
			logger.L.Panic().Err(err).Msgf("Failed to parse font %s", fp)
		}
		large := truetype.NewFace(fnt, &truetype.Options{Size: 40})
		small := truetype.NewFace(fnt, &truetype.Options{Size: 25})
		ass.Fonts[noExt+"Large"] = large
		ass.Fonts[noExt+"Small"] = small

		logger.L.Info().Msgf("Loaded %s", fp)
	}
	return &ass
}
