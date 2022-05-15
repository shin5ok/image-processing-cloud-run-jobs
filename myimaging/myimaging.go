package myimaging

import (
	"fmt"

	"github.com/disintegration/imaging"
)

type Image struct {
	Filename string
}

func (i Image) MakeSmall(width int) (string, error) {
	newfilename := fmt.Sprintf("%s_%s", "new", i.Filename)
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
