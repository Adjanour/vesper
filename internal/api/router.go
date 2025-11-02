package api

import (
	"encoding/json"
	"net/http"

	"github.com/Adjanour/vesper/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type APIRouter struct {
	router *chi.Mux
	db     *database.Queries
}

func NewAPIRouter(q *database.Queries) *chi.Mux {
	api := &APIRouter{
		router: chi.NewRouter(),
		db:     q,
	}
	return api.Routes()
}



// borrowed from github.com/AmoabaKelvin/loglevel/internal/api
func WriteJsonResponse(w http.ResponseWriter, status int, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(payload)
}

func (ar *APIRouter) Routes() *chi.Mux {
	ar.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	ar.router.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			WriteJsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
		})
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/{id}", ar.getTask)
			r.Post("/", ar.createTask)
			r.Delete("/{id}", ar.deleteTask)
		})
	})

	return ar.router
}
