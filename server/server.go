package server

import (
	"github.com/SigurdRiseth/CountryInfoService/config"
	"github.com/SigurdRiseth/CountryInfoService/handlers"
	"github.com/SigurdRiseth/CountryInfoService/utils"
	"log"
	"net/http"
	"time"
)

// StartServer initializes and starts the HTTP server
func StartServer() {
	handler.StartTime = time.Now() // Initialize start time

	// Load environment variables
	config.LoadEnvVariables()

	// Get the port from environment variables, default to 8080
	port := config.GetPort()

	// Instantiate the router
	router := setupRouter()

	// Start the server
	log.Println("Server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router)) // TODO: Gracefully shutdown server?
}

// setupRouter sets up the HTTP routes and handlers
func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Define the endpoints
	router.HandleFunc(utils.INFO_PATH, makeHTTPHandleFunc(handler.HandleInfo))
	router.HandleFunc(utils.POPULATION_PATH, makeHTTPHandleFunc(handler.HandlePopulation))
	router.HandleFunc(utils.STATUS_PATH, makeHTTPHandleFunc(handler.HandleStatus))
	router.HandleFunc("/", makeHTTPHandleFunc(handler.DefaultHandler))

	return router
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Error string `json:"message"`
}

// makeHTTPHandleFunc is a helper function that wraps an apiFunc with error handling.
// It returns an http.HandlerFunc that logs the error and sends an HTTP 500 status code
// if the wrapped function returns an error.
//
// Parameters:
// - f: The apiFunc to be wrapped.
//
// Returns:
// - http.HandlerFunc: A function that handles HTTP requests and responses.
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// TODO: switch on error type and return appropriate status code
		}
	}
}
