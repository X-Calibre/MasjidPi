package main

import (
	"fmt"
	"log"
	"os"

	"github.com/X-Calibre/MasjidPi/backend/internal/app"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

func main() {
	if handleCommandLine(os.Args[1:]) {
		return
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func handleCommandLine(args []string) bool {
	if len(args) == 0 {
		return false
	}

	switch args[0] {
	case "--help", "-h":
		printUsage()
		return true
	case "--version", "-v":
		fmt.Println(version.Version)
		return true
	default:
		return false
	}
}

func printUsage() {
	fmt.Printf(`%s - lightweight internet radio for masājid streams

Usage:
  masjidpi [options]

Options:
  -h, --help       Show this help message and exit
  -v, --version    Show the application version and exit

With no options, MasjidPi starts the application and HTTP server.
`, version.AppName)
}
