package main

import (
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/app"
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
