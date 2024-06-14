package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// define port constant value
const port int = 8080

func main() {
	// Create chi router
	router := chi.NewRouter()
	// Configure chi logging
	router.Use(middleware.Logger)
	// Create router
	router.Get("/hello", basicHandler)

	// create http server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port), // convert port to ASCII
		Handler: router,
	}

	// configure logs to be written to stdout
	log.SetOutput(log.Writer())

	// Start the server
	log.Printf("server is listening on port %d", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("server creation error occured:", err)
	}
}

// first handler
func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go!"))
}