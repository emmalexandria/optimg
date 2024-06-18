package main

import (
	"math"

	"github.com/h2non/bimg"
)

func calcSteps(maxWidth int, minWidth int) int {
	return int(math.Round(float64(maxWidth) / float64(minWidth)))
}

func convFormat(image bimg.Image, format bimg.ImageType) (bimg.Image, error) {
	newBuf, err := image.Convert(format)
	if err != nil {
		return image, err
	}

	newImg := bimg.NewImage(newBuf)
	return *newImg, nil
}

func resizeRatio(image bimg.Image, ratio float64) (bimg.Image, error) {
	size, err := image.Size()
	if err != nil {
		return image, err
	}

	newBuf, err := image.Resize(int(float64(size.Width)*ratio), int(float64(size.Height)*ratio))
	if err != nil {
		return image, err
	}

	newImg := bimg.NewImage(newBuf)
	return *newImg, nil
}

func processImage(image bimg.Image, options bimg.Options) (bimg.Image, error) {
	newBuf, err := image.Process(options)
	if err != nil {
		return image, err
	}

	return *bimg.NewImage(newBuf), nil
}

/* func getResizedImages(image bimg.Image, steps int) ([]bimg.Image, error) {
	ratio := 1.0 / steps
	images := []bimg.Image{image}

	for i := range steps - 1 {

	}
} */

func readImage(path string) (bimg.Image, error) {
	buffer, err := bimg.Read(path)
	if err != nil {
		return *bimg.NewImage([]byte{}), err
	}

	image := bimg.NewImage(buffer)
	return *image, nil

}
