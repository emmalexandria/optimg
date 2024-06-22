package main

import (
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
