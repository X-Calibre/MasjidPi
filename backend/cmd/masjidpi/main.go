package main

import (
	"log"

	"github.com/X-Calibre/MasjidPi/backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
