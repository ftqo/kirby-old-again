package files

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strings"
)

type Hoarder struct {
	Fonts  map[string][]byte
	Images map[string]image.Image
}

func (h *Hoarder) LoadFiles() {
	h.Images = make(map[string]image.Image)
	h.Fonts = make(map[string][]byte)
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	assetsPath := path.Join(d, "../assets")
	imagesPath := path.Join(assetsPath, "images")
	fontsPath := path.Join(assetsPath, "fonts")
	imgs, err := ioutil.ReadDir(imagesPath)
	if err != nil {
		log.Panicln(err)
	}
	for _, file := range imgs {
		fp := path.Join(imagesPath, file.Name())
		bts, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Panicln(err)
		}
		fn := file.Name()
		noPre := fn[strings.LastIndex(file.Name(), "-")+1:]
		noExt := noPre[:strings.Index(noPre, ".")]
		h.Images[noExt], _, err = image.Decode(bytes.NewReader(bts))
		if err != nil {
			log.Panicf("Error decoding %s: %v", fp, err)
		}
		log.Printf("Loaded %s", fp)
	}
	fonts, err := ioutil.ReadDir(fontsPath)
	if err != nil {
		log.Panicln(err)
	}
	for _, file := range fonts {
		fp := path.Join(fontsPath, file.Name())
		bts, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Panicln(err)
		}
		fn := file.Name()
		noPre := fn[strings.LastIndex(file.Name(), "-")+1:]
		noExt := noPre[:strings.Index(noPre, ".")]
		h.Fonts[noExt] = bts
		log.Printf("Loaded %s", fp)
	}
}
