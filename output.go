package main

import (
	"fmt"
	"github.com/fatih/color"
)

func successPrintln(output string, params ...interface{}) {
	color.New(color.FgGreen, color.Bold).Printf(output+"\n", params...)
}

func errorPrintln(output string, params ...interface{}) {
	color.New(color.FgRed, color.Bold).Printf(output+"\n", params...)
}

func warningPrintln(output string, params ...interface{}) {
	color.New(color.FgYellow).Printf(output+"\n", params...)
}

func infoPrintln(output string, params ...interface{}) {
	color.New(color.FgWhite).Printf(output+"\n", params...)
}

func alreadyProcessedMessage(file string) {
	warningPrintln("Already processed %s, skipping...", file)
}

func generateSrcSet(images []imagePath) (string, error) {
	boldColor := color.New(color.FgHiWhite, color.Bold)
	srcSetStr := boldColor.Sprintf("%s:", images[0].originalPath)
	for _, image := range images {
		size, err := image.image.Size()
		if err != nil {
			return "", err
		}
		srcSetStr += fmt.Sprintf("%s %v%s ", image.path, size.Width, "w")
	}
	return srcSetStr, nil
}
