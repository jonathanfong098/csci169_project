package main

import (
	"log"
	"net/http" // Import the net/http package to use http.Server
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors" // Import the cors package
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set.")
	}

	// Initialize the router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	router.Mount("/v1", v1Router)

	// Define the server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port, // Make sure this line has a colon before the port and a comma at the end
	}

	log.Printf("Server is running on port %s", port)

	// Start listening and serving HTTP requests
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
