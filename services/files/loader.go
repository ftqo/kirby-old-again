package files

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type HoardedFiles struct {
	Images []image.Image
}

func (f *HoardedFiles) LoadImages() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	assets := path.Join(d, "../../assets/images")
	var filePaths []string
	err := filepath.Walk(assets, visit(&filePaths))
	filePaths = filePaths[1:]
	if err != nil {
		log.Panicln(err)
	}
	for _, fp := range filePaths {
		fmt.Println(fp)
		file, err := os.ReadFile(fp)
		if err != nil {
			log.Panicln(err)
		}
		img, _, err := image.Decode(bytes.NewReader(file))
		if err != nil {
			log.Panicln(err)
		}
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
