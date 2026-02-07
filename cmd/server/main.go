package main

import (
	"log"
	"net/http"

	"github.com/Adjanour/vesper/internal/api"
	"github.com/Adjanour/vesper/internal/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queries := database.NewQueries(db)
	apiRouter := api.NewAPIRouter(queries)

	// Create main router
	mainRouter := chi.NewRouter()
	
	// Mount API routes (apiRouter already has /api prefix in its routes)
	mainRouter.Mount("/", apiRouter)
	
	// Serve index.html for root path
	mainRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	log.Println("Server starting on :8080")
	log.Println("Web UI available at http://localhost:8080")
	log.Println("API available at http://localhost:8080/api")
	http.ListenAndServe(":8080", mainRouter)
}
