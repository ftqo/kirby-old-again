package files

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Hoarder struct {
	Fonts  map[string]font.Face
	Images map[string]image.Image
}

func (h *Hoarder) LoadFiles() {
	h.Images = make(map[string]image.Image)
	h.Fonts = make(map[string]font.Face)
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	assetsPath := path.Join(d, "../assets")
	imagesPath := path.Join(assetsPath, "images")
	fontsPath := path.Join(assetsPath, "fonts")
	imgs, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		log.Panicf("failed to read directory %s: %v", imgs, err)
	}
	for _, file := range imgs {
		fp := path.Join(imagesPath, file.Name())
		bts, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Panicf("failed to read file %s: %v", fp, err)
		}
		fn := file.Name()
		noPre := fn[strings.LastIndex(file.Name(), "-")+1:]
		noExt := noPre[:strings.Index(noPre, ".")]
		h.Images[noExt], _, err = image.Decode(bytes.NewReader(bts))
		if err != nil {
			log.Panicf("failed to decode %s: %v", fp, err)
		}
		log.Printf("loaded %s !", fp)
	}
	fonts, err := ioutil.ReadDir(fontsPath)
	if err != nil {
		log.Panicf("failed to read directory %s: %v", fonts, err)
	}
	for _, file := range fonts {
		fp := path.Join(fontsPath, file.Name())
		bts, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Panicf("failed to read file %s: %v", fp, err)
		}
		fn := file.Name()
		noPre := fn[strings.LastIndex(file.Name(), "-")+1:]
		noExt := noPre[:strings.Index(noPre, ".")]
		font, err := truetype.Parse(bts)
		if err != nil {
			log.Panicf("failed to parse font %s: %v", fp, err)
		}
		large := truetype.NewFace(font, &truetype.Options{Size: 40})
		small := truetype.NewFace(font, &truetype.Options{Size: 25})
		h.Fonts[noExt+"Large"] = large
		h.Fonts[noExt+"Small"] = small

		log.Printf("loaded %s !", fp)
	}
}
