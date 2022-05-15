package myimaging

import (
	"fmt"

	"github.com/disintegration/imaging"
)

type Image struct {
	Filename string
}

type imageInterface interface {
	MakeSmall(filename string) (string, error)
}

func (i Image) MakeSmall(filename string) (string, error) {
	newfilename := fmt.Sprintf("%s_%s", "new", i.Filename)
	srcImage, err := imaging.Open(filename, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	newfilenameImage := imaging.Resize(srcImage, 800, 0, imaging.Lanczos)
	if err := imaging.Save(newfilenameImage, newfilename); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return newfilename, nil
}
