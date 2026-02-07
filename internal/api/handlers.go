package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adjanour/vesper/internal/database"
	"github.com/Adjanour/vesper/internal/models"
	"github.com/go-chi/chi/v5"
)

// validateTask validates task fields
func validateTask(t *models.Task) error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if t.UserID == "" {
		return errors.New("user_id is required")
	}
	if t.Start.IsZero() {
		return errors.New("start time is required")
	}
	if t.End.IsZero() {
		return errors.New("end time is required")
	}
	if t.End.Before(t.Start) || t.End.Equal(t.Start) {
		return errors.New("end time must be after start time")
	}
	if !models.IsValidStatus(t.Status) {
		return errors.New("invalid status")
	}
	return nil
}

func (ar *APIRouter) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromRequest(r)
	tasks, err := ar.db.GetTasks(ctx, userID)
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

	if userID, ok := userIDFromHeader(r); ok && task.UserID != userID {
		http.Error(w, "task not found", http.StatusNotFound)
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

	if userID, ok := userIDFromHeader(r); ok {
		t.UserID = userID
	}

	// Validate task
	if err := validateTask(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set default status if not provided
	if t.Status == "" {
		t.Status = models.StatusScheduled
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

func (ar *APIRouter) updateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Set ID from URL parameter
	t.ID = id

	if userID, ok := userIDFromHeader(r); ok {
		existing, err := ar.db.GetTask(ctx, id)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				http.Error(w, "task not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if existing.UserID != userID {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		t.UserID = userID
	}

	// Validate task
	if err := validateTask(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ar.db.UpdateTask(ctx, t); err != nil {
		switch err {
		case database.ErrNotFound:
			http.Error(w, "task not found", http.StatusNotFound)
		case database.ErrTaskOverlap:
			http.Error(w, "task overlaps with existing task", http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	WriteJsonResponse(w, http.StatusOK, t)
}

func (ar *APIRouter) deleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if userID, ok := userIDFromHeader(r); ok {
		task, err := ar.db.GetTask(ctx, id)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				http.Error(w, "task not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if task.UserID != userID {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
	}

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

func userIDFromRequest(r *http.Request) string {
	if userID, ok := userIDFromHeader(r); ok {
		return userID
	}
	return "1"
}

func userIDFromHeader(r *http.Request) (string, bool) {
	userID := r.Header.Get("X-User-ID")
	return userID, userID != ""
}
