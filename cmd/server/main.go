package main

import (
	"log"

	app "github.com/RubenRodrigo/go-tiny-store/internal"
)

func main() {
	application := app.New()

	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

}
