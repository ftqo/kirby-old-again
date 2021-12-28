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
	Images []image.Image
}

func (f *Hoarder) LoadImages() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	assets := path.Join(d, "../assets/images")
	var filePaths []string
	err := filepath.Walk(assets, visit(&filePaths))
	filePaths = filePaths[1:]
	if err != nil {
		log.Panicln(err)
	}
	for _, fp := range filePaths {
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
