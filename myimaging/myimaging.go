package myimaging

import (
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type Image struct {
	Filename string
}

func (i Image) MakeSmall(width int) (string, error) {

	dirname := filepath.Dir(i.Filename)
	basename := filepath.Base(i.Filename)
	newFilename := fmt.Sprintf("%s/%s_%s", dirname, "new", basename)

	srcImage, err := imaging.Open(i.Filename, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	newFilenameImage := imaging.Resize(srcImage, width, 0, imaging.Lanczos)
	if err := imaging.Save(newFilenameImage, newFilename); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return newFilename, nil
}
