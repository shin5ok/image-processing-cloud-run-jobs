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
	newfilename := fmt.Sprintf("%s/%s_%s", dirname, "new", basename)

	srcImage, err := imaging.Open(i.Filename, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	newfilenameImage := imaging.Resize(srcImage, width, 0, imaging.Lanczos)
	if err := imaging.Save(newfilenameImage, newfilename); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return newfilename, nil
}
