package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adjanour/vesper/internal/database"
	"github.com/Adjanour/vesper/internal/models"
	"github.com/go-chi/chi/v5"
)

func (ar *APIRouter) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := ar.db.GetTasks(ctx,"1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(w, http.StatusOK, map[string]any{
		"tasks": tasks,
	})
}

func (ar *APIRouter) getTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	task, err := ar.db.GetTask(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJsonResponse(w, http.StatusOK, task)
}

func (ar *APIRouter) createTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := ar.db.CreateTask(ctx, t); err != nil {
		switch err {
		case database.ErrTaskOverlap:
			http.Error(w, "task overlaps with existing task", http.StatusConflict)
		case database.ErrInvalid:
			http.Error(w, "invalid user id", http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	WriteJsonResponse(w, http.StatusCreated, t)
}

func (ar *APIRouter) deleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := ar.db.DeleteTask(ctx, id); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
