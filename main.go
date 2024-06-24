package main

import (
	"fmt"
	"path"
	"slices"

	"github.com/alexflint/go-arg"
	"github.com/h2non/bimg"
	"github.com/schollz/progressbar/v3"
)

type Args struct {
	Input           []string `arg:"required, positional" help:"input directories or files"`
	InputExtensions []string `arg:"-i, --input" `
	Steps           int      `default:"4" arg:"-s, --steps" help:"resizing steps"`
	Quality         int      `arg:"-q, --quality" default:"80" help:"output quality"`
	OutputDir       string   `default:"processed" arg:"-o, --output" help:"directory to output to relative to input"`
	OutputType      string   `arg:"-t, --type" default:".webp"`

	Recursive      bool `arg:"-r, --recurse" help:"recurse through directories"`
	ClearOutputDir bool `arg:"-c, --clear" help:"delete output directory before writing"`
	StripMetadata  bool `arg:"-s, --strip" help:"strip file metadata"`
	GenerateSrcSet bool `arg:"-g, --genSrcSet"`
}

func (Args) Version() string {
	return "optimg 1.0.0"
}

var Progress struct {
	bar *progressbar.ProgressBar
}

func main() {
	var args Args
	arg.MustParse(&args)

	args.InputExtensions = []string{".webp", ".png", ".jpg", ".jpeg"}

	options := bimg.Options{StripMetadata: args.StripMetadata, Quality: args.Quality}

	count := 0
	for _, input := range args.Input {
		estimate, err := estTotal(args, input)
		if err != nil {
			return
		}
		count += estimate
	}

	Progress.bar = progressbar.NewOptions(count, progressbar.OptionShowCount(), progressbar.OptionShowIts(), progressbar.OptionFullWidth(), progressbar.OptionSetTheme(progressbar.Theme{
		Saucer:        "█",
		SaucerHead:    "█",
		BarStart:      "[",
		BarEnd:        "]",
		SaucerPadding: "▒",
	}))

	alreadyProcessed := []string{}
	srcSetList := []string{}
	for _, inputPath := range args.Input {
		processed, srcSet, err := processInput(args, inputPath, options, alreadyProcessed)

		alreadyProcessed = processed
		if err != nil {
			continue
		}

		srcSetList = append(srcSetList, srcSet...)
	}

	Progress.bar.Finish()
	if args.GenerateSrcSet {
		for _, srcSet := range srcSetList {
			fmt.Println(srcSet)
		}
	}

}

func processInput(args Args, inputPath string, options bimg.Options, alreadyProcessed []string) ([]string, []string, error) {
	isDir, err := isDir(inputPath)
	if err != nil {
		return alreadyProcessed, nil, nil
	}
	srcSetList := []string{}
	if isDir {
		dirs := []string{}
		if args.Recursive {
			subDirs, err := getSubDirs(inputPath, args.OutputDir)
			if err != nil {
				errorPrintln("Error getting subdirectories of %s", inputPath)
				return alreadyProcessed, nil, err
			}
			dirs = append(dirs, subDirs...)
		} else {
			dirs = append(dirs, inputPath)
		}
		processed, srcSet, err := processDirList(args, dirs, alreadyProcessed, options)
		alreadyProcessed = processed
		if err != nil {
			return alreadyProcessed, nil, err
		}
		srcSetList = append(srcSetList, srcSet...)

	} else {
		dir := path.Dir(inputPath)
		err = initOutputDir(dir, args.OutputDir, args.ClearOutputDir, alreadyProcessed)
		if err != nil {
			errorPrintln("An error occured creating the output dir: %s", err.Error())
			return alreadyProcessed, nil, err
		}

		processed, srcSet, err := processFile(args, inputPath, dir, alreadyProcessed, options)
		alreadyProcessed = processed
		if err != nil {
			errorPrintln("An error occured processing %s: %s", inputPath, err.Error())
			return alreadyProcessed, nil, err
		}
		srcSetList = append(srcSetList, srcSet)
	}

	return alreadyProcessed, srcSetList, nil
}

func processDirList(args Args, dirs []string, alreadyProcessed []string, imageOptions bimg.Options) ([]string, []string, error) {
	srcSetList := []string{}
	for _, dir := range dirs {
		images, err := getFilesInDir(dir, args.InputExtensions)
		if err != nil {
			return alreadyProcessed, nil, err
		}

		err = initOutputDir(dir, args.OutputDir, args.ClearOutputDir, alreadyProcessed)
		if err != nil {
			errorPrintln("An error occured creating the output dir: %s", err.Error())
			return alreadyProcessed, nil, err
		}
		for _, img := range images {
			processed, srcSet, err := processFile(args, img, dir, alreadyProcessed, imageOptions)
			alreadyProcessed = processed
			if err != nil {
				errorPrintln("An error occured processing %s: %s", img, err.Error())
				return alreadyProcessed, nil, err
			}
			srcSetList = append(srcSetList, srcSet)
		}

	}

	return alreadyProcessed, srcSetList, nil
}

func processFile(args Args, filePath string, dir string, alreadyProcessed []string, imageOptions bimg.Options) ([]string, string, error) {
	if slices.Contains(alreadyProcessed, filePath) {
		return alreadyProcessed, "", nil
	}
	Progress.bar.Describe(fmt.Sprintf("Processing %s", path.Base(filePath)))
	resizedImages, err := processImage(filePath, imageOptions, args.Steps)
	if err != nil {
		errorPrintln("Error resizing image %s: %s", filePath, err.Error())
		return alreadyProcessed, "", err
	}

	imagePaths, err := saveImagesInOutput(resizedImages, filePath, args.OutputType, args.OutputDir)
	if err != nil {
		errorPrintln("Error saving resized images for %s: %s", err.Error())
		return alreadyProcessed, "", err
	}
	Progress.bar.Add(1)

	srcSet := ""
	if args.GenerateSrcSet {
		srcSet, err = generateSrcSet(imagePaths)
		if err != nil {
			return nil, "", err
		}
	}

	return append(alreadyProcessed, filePath), srcSet, nil
}
