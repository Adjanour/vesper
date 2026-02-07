package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Adjanour/vesper/internal/database"
	"github.com/Adjanour/vesper/internal/models"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *database.Queries {
	// Use in-memory SQLite for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			start DATETIME NOT NULL,
			end DATETIME NOT NULL,
			status TEXT NOT NULL CHECK (status IN ('scheduled', 'deleted', 'replaced')),
			user_id TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_start_end ON tasks(start, end);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Insert test user
	_, err = db.ExecContext(context.Background(), "INSERT OR IGNORE INTO users (id, username) VALUES ('1', 'testuser')")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	return database.NewQueries(db)
}

func TestHealthEndpoint(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestCreateTask(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	task := models.Task{
		ID:     "test-create-001",
		Title:  "Test Task",
		Start:  time.Now().Add(1 * time.Hour),
		End:    time.Now().Add(2 * time.Hour),
		UserID: "test-user",
		Status: models.StatusScheduled,
	}

	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response models.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Title != task.Title {
		t.Errorf("Expected title '%s', got '%s'", task.Title, response.Title)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	tests := []struct {
		name       string
		task       map[string]interface{}
		wantStatus int
		wantError  string
	}{
		{
			name: "missing title",
			task: map[string]interface{}{
				"id":      "test-001",
				"start":   time.Now().Format(time.RFC3339),
				"end":     time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				"user_id": "test-user",
				"status":  "scheduled",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "title is required",
		},
		{
			name: "missing user_id",
			task: map[string]interface{}{
				"id":     "test-002",
				"title":  "Test",
				"start":  time.Now().Format(time.RFC3339),
				"end":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				"status": "scheduled",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "user_id is required",
		},
		{
			name: "invalid status",
			task: map[string]interface{}{
				"id":      "test-003",
				"title":   "Test",
				"start":   time.Now().Format(time.RFC3339),
				"end":     time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				"user_id": "test-user",
				"status":  "invalid",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid status",
		},
		{
			name: "end before start",
			task: map[string]interface{}{
				"id":      "test-004",
				"title":   "Test",
				"start":   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
				"end":     time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				"user_id": "test-user",
				"status":  "scheduled",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "end time must be after start time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.task)
			req := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if w.Body.String() != tt.wantError+"\n" {
				t.Errorf("Expected error '%s', got '%s'", tt.wantError, w.Body.String())
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	// Create a task first
	task := models.Task{
		ID:     "test-get-001",
		Title:  "Test Get Task",
		Start:  time.Now().Add(1 * time.Hour),
		End:    time.Now().Add(2 * time.Hour),
		UserID: "test-user",
		Status: models.StatusScheduled,
	}

	body, _ := json.Marshal(task)
	createReq := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	// Now get the task
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/test-get-001", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != task.ID {
		t.Errorf("Expected ID '%s', got '%s'", task.ID, response.ID)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestListTasks(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	// Create two tasks
	tasks := []models.Task{
		{
			ID:     "test-list-001",
			Title:  "Task 1",
			Start:  time.Now().Add(1 * time.Hour),
			End:    time.Now().Add(2 * time.Hour),
			UserID: "1",
			Status: models.StatusScheduled,
		},
		{
			ID:     "test-list-002",
			Title:  "Task 2",
			Start:  time.Now().Add(3 * time.Hour),
			End:    time.Now().Add(4 * time.Hour),
			UserID: "1",
			Status: models.StatusScheduled,
		},
	}

	for _, task := range tasks {
		body, _ := json.Marshal(task)
		createReq := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createReq)
	}

	// List tasks
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string][]*models.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response["tasks"]) < 2 {
		t.Errorf("Expected at least 2 tasks, got %d", len(response["tasks"]))
	}
}

func TestUpdateTask(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	// Create a task first
	task := models.Task{
		ID:     "test-update-001",
		Title:  "Original Title",
		Start:  time.Now().Add(1 * time.Hour),
		End:    time.Now().Add(2 * time.Hour),
		UserID: "test-user",
		Status: models.StatusScheduled,
	}

	body, _ := json.Marshal(task)
	createReq := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	// Update the task
	updatedTask := task
	updatedTask.Title = "Updated Title"

	updateBody, _ := json.Marshal(updatedTask)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/tasks/test-update-001", bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()

	router.ServeHTTP(updateW, updateReq)

	if updateW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", updateW.Code, updateW.Body.String())
	}

	var response models.Task
	if err := json.NewDecoder(updateW.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", response.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	// Create a task first
	task := models.Task{
		ID:     "test-delete-001",
		Title:  "To Be Deleted",
		Start:  time.Now().Add(1 * time.Hour),
		End:    time.Now().Add(2 * time.Hour),
		UserID: "test-user",
		Status: models.StatusScheduled,
	}

	body, _ := json.Marshal(task)
	createReq := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	// Delete the task
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/tasks/test-delete-001", nil)
	deleteW := httptest.NewRecorder()

	router.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", deleteW.Code)
	}

	// Verify it's deleted
	getReq := httptest.NewRequest(http.MethodGet, "/api/tasks/test-delete-001", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 after deletion, got %d", getW.Code)
	}
}

func TestTaskOverlapDetection(t *testing.T) {
	queries := setupTestDB(t)
	router := NewAPIRouter(queries)

	// Create first task
	task1 := models.Task{
		ID:     "test-overlap-001",
		Title:  "First Task",
		Start:  time.Date(2026, 2, 8, 9, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC),
		UserID: "test-user",
		Status: models.StatusScheduled,
	}

	body1, _ := json.Marshal(task1)
	createReq1 := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body1))
	createReq1.Header.Set("Content-Type", "application/json")
	createW1 := httptest.NewRecorder()
	router.ServeHTTP(createW1, createReq1)

	// Try to create overlapping task
	task2 := models.Task{
		ID:     "test-overlap-002",
		Title:  "Overlapping Task",
		Start:  time.Date(2026, 2, 8, 9, 30, 0, 0, time.UTC),
		End:    time.Date(2026, 2, 8, 10, 30, 0, 0, time.UTC),
		UserID: "test-user",
		Status: models.StatusScheduled,
	}

	body2, _ := json.Marshal(task2)
	createReq2 := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body2))
	createReq2.Header.Set("Content-Type", "application/json")
	createW2 := httptest.NewRecorder()
	router.ServeHTTP(createW2, createReq2)

	if createW2.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for overlapping task, got %d", createW2.Code)
	}
}
