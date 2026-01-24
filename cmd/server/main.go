package main

import (
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/app"
)

func main() {
	// router := mux.NewRouter()
	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("INSIDE API IN /")
	// 	w.Write([]byte("Hello World! OUTSIDE"))
	// }).Methods("GET")

	// addr := fmt.Sprintf("localhost:3001")
	// http.ListenAndServe(addr, router)

	application := app.New()

	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

}
