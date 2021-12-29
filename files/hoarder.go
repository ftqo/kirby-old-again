package files

import (
	"bytes"
	"image"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type Hoarder struct {
	Font   []byte
	Images []image.Image
}

func (f *Hoarder) LoadFiles() {
	// Loading images
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	assetsPath := path.Join(d, "../assets")
	imagesPath := path.Join(assetsPath, "/images")
	var imagePaths []string
	err := filepath.Walk(imagesPath, visit(&imagePaths))
	if err != nil {
		log.Panicln(err)
	}
	imagePaths = imagePaths[1:]
	for _, fp := range imagePaths {
		file, err := os.ReadFile(fp)
		if err != nil {
			log.Panicln(err)
		}
		img, _, err := image.Decode(bytes.NewReader(file))
		if err != nil {
			log.Panicln(err)
		}
		log.Printf("Loaded %s", fp)
		f.Images = append(f.Images, img)
	}

	// Loading font
	fontsPath := path.Join(assetsPath, "/fonts")
	var fontPaths []string
	err = filepath.Walk(fontsPath, visit(&fontPaths))
	if err != nil {
		log.Panicln(err)
	}
	fontPaths = fontPaths[1:]
	file, err := os.ReadFile(fontPaths[0])
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("Loaded %s", fontPaths[0])
	f.Font = file
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}
