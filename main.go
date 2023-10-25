package main

import (
	"cognivaultServer/api"
	"cognivaultServer/database"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to the SQLite database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables for collections, tags, and data points
	err = database.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	// Set up the chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Set up the API routes
	api.SetRoutes(r)

	// Serve the Swagger UI for API documentation
	r.Get("/swagger/*", api.SwaggerHandler())

	// Start the server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
