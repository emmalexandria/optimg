package main

import (
	"path"
	"slices"

	"github.com/alexflint/go-arg"
	"github.com/h2non/bimg"
)

type Args struct {
	Input           []string `arg:"required, positional" help:"input directories or files"`
	InputExtensions []string `arg:"-i, --input" `
	Steps           int      `default:"4" arg:"-s, --steps" help:"resizing steps"`
	Quality         int      `arg:"-q, --quality" default:"80" help:"output quality"`
	OutputDir       string   `default:"processed" arg:"-o, --output" help:"directory to output to relative to input"`
	OutputType      string   `arg:"-t, --type" default:".webp"`
	IgnoreDir       string   `arg:"-i, --ignore" help:"directory to ignore. same as output by default"`

	Recursive      bool `arg:"-r, --recurse" help:"recurse through directories"`
	ClearOutputDir bool `arg:"-c, --clear" help:"delete output directory before writing"`
	StripMetadata  bool `arg:"-s, --strip" help:"strip file metadata"`
}

func (Args) Version() string {
	return "optimg 1.0.0"
}

func main() {
	var args Args
	args.InputExtensions = []string{".webp", ".png", ".jpg", ".jpeg"}
	arg.MustParse(&args)

	if args.IgnoreDir == "" {
		args.IgnoreDir = args.OutputDir
	}
	options := bimg.Options{StripMetadata: args.StripMetadata, Quality: args.Quality}

	alreadyProcessed := []string{}
	for _, inputPath := range args.Input {
		processed, err := processInput(args, inputPath, args.InputExtensions, options, alreadyProcessed)

		alreadyProcessed = processed
		if err != nil {
			continue
		}
	}
}

func processInput(args Args, inputPath string, extensions []string, options bimg.Options, alreadyProcessed []string) ([]string, error) {
	isDir, err := isDir(inputPath)
	if err != nil {
		return alreadyProcessed, nil
	}

	if isDir {
		dirs := []string{}
		if args.Recursive {
			subDirs, err := getSubDirs(inputPath, args.OutputDir)
			if err != nil {
				errorPrintln("Error getting subdirectories of %s", inputPath)
				return alreadyProcessed, err
			}
			dirs = append(dirs, subDirs...)
		} else {
			dirs = append(dirs, inputPath)
		}
		processed, err := processDirList(args, dirs, extensions, alreadyProcessed, options)
		alreadyProcessed = processed
		if err != nil {
			return alreadyProcessed, err
		}

	} else {
		dir := path.Dir(inputPath)
		err = initOutputDir(dir, args.OutputDir, args.ClearOutputDir, alreadyProcessed)
		if err != nil {
			errorPrintln("An error occured creating the output dir: %s", err.Error())
			return alreadyProcessed, err
		}

		processed, err := processFile(args, inputPath, dir, extensions, alreadyProcessed, options)
		alreadyProcessed = processed
		if err != nil {
			errorPrintln("An error occured processing %s: %s", inputPath, err.Error())
			return alreadyProcessed, err
		}
	}

	return alreadyProcessed, nil
}

func processDirList(args Args, dirs []string, extensions []string, alreadyProcessed []string, imageOptions bimg.Options) ([]string, error) {
	for _, dir := range dirs {
		images, err := getFilesInDir(dir, extensions)
		if err != nil {
			return alreadyProcessed, nil
		}

		err = initOutputDir(dir, args.OutputDir, args.ClearOutputDir, alreadyProcessed)
		if err != nil {
			errorPrintln("An error occured creating the output dir: %s", err.Error())
			return alreadyProcessed, err
		}
		for _, img := range images {
			processed, err := processFile(args, img, dir, extensions, alreadyProcessed, imageOptions)
			alreadyProcessed = processed
			if err != nil {
				errorPrintln("An error occured processing %s: %s", img, err.Error())
				return alreadyProcessed, err
			}
		}

	}

	return alreadyProcessed, nil
}

func processFile(args Args, filePath string, dir string, extensions []string, alreadyProcessed []string, imageOptions bimg.Options) ([]string, error) {
	if slices.Contains(alreadyProcessed, filePath) {
		return alreadyProcessed, nil
	}
	resizedImages, err := processImage(filePath, imageOptions, args.Steps)
	if err != nil {
		errorPrintln("Error resizing image %s: %s", filePath, err.Error())
		return alreadyProcessed, err
	}

	err = saveImagesInOutput(resizedImages, filePath, args.OutputType, args.OutputDir)
	if err != nil {
		errorPrintln("Error saving resized images for %s: %s", err.Error())
		return alreadyProcessed, err
	}
	successPrintln("Sucessfully processed and saved %s", filePath)

	return append(alreadyProcessed, filePath), nil

}
