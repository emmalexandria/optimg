package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/h2non/bimg"
)

type args struct {
	Input      []string `arg:"required, positional" help:"input directories or files"`
	Steps      int      `default:"4" arg:"-s, --steps" help:"resizing steps. takes priority over min width"`
	MinWidth   int      `arg:"-m, --min_width" help:"desired width of smallest image"`
	Quality    int      `arg:"-q, --quality" default:"80" help:"output quality"`
	OutputDir  string   `default:"processed" arg:"-o, --output" help:"directory to output to relative to input"`
	OutputType string   `arg:"-t, --type" default:".webp"`
	IgnoreDir  string   `arg:"-i, --ignore" help:"directory to ignore. same as output by default"`

	Recursive      bool `arg:"-r, --recurse"`
	ClearOutputDir bool `arg:"-c, --clear"`
	Verbose        bool `arg:"-v, --verbose"`
	StripMetadata  bool `arg:"-s, --strip"`
}

func (args) Version() string {
	return "optimg 1.0.0"
}

func main() {
	var args args
	arg.MustParse(&args)

	if args.IgnoreDir == "" {
		args.IgnoreDir = args.OutputDir
	}

	inputTypes := []string{".webp", ".png", ".jpeg", ".jpg"}
	options := bimg.Options{StripMetadata: args.StripMetadata, Quality: args.Quality}

	fmt.Printf("%+v", args)

	imagePaths := []string{}
	for _, path := range args.Input {
		images, err := getInputFiles(path, inputTypes, args.Recursive)
		if err != nil {
			fmt.Printf("Failed to read input path %s \n", path)
		}

		imagePaths = append(imagePaths, images...)
	}

	for _, image := range imagePaths {
		imageData, err := readImage(image)
		if err != nil {
			fmt.Printf("Failed to read image %s \n", image)
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		outputImages := []bimg.Image{}

		imageData, err = convFormat(imageData, bimg.WEBP)
		imageData, err = processImage(imageData, options)

		if err != nil {
			fmt.Printf("Failed to convert image %s \n", image)
			fmt.Fprintln(os.Stderr, err)
			continue

		}
		if err != nil {
			fmt.Printf("Failed to convert image %s \n", image)
			fmt.Fprintln(os.Stderr, err)
			continue

		}
		for i := args.Steps; i > 0; i-- {
			ratio := float64(i) / float64(args.Steps)
			resizedImage, err := resizeRatio(imageData, ratio)
			if err != nil {
				fmt.Printf("Failed to resize image %s \n", image)
			}

			outputImages = append(outputImages, resizedImage)
		}

		for _, img := range outputImages {
			ext := filepath.Ext(image)

			imgName, _ := getModifiedImageName(image[0:len(image)-len(ext)], args.OutputType, img, 0)
			fmt.Println(imgName)

		}
		//fmt.Printf("Created %v output images for %s", len(outputImages), image)

	}

}
