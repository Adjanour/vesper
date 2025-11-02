package main

import (
	"log"
	"net/http"

	"github.com/Adjanour/vesper/internal/api"
	"github.com/Adjanour/vesper/internal/database"
)

func main() {
	db,err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queries := database.NewQueries(db)
	apiRouter := api.NewAPIRouter(queries)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", apiRouter)
}
