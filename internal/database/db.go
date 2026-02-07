package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Adjanour/vesper/internal/models"
	_ "modernc.org/sqlite"
)

// Domain-level errors
var (
	ErrNotFound     = errors.New("not found")
	ErrDuplicate    = errors.New("duplicate")
	ErrInvalid      = errors.New("invalid")
	ErrUnauthorized = errors.New("unauthorized")
	ErrTaskOverlap  = errors.New("task overlap")
)

// DBTX interface allows mocking or using transactions
type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// Queries struct wraps DBTX (could be *sql.DB or *sql.Tx)
type Queries struct {
	db DBTX
}

func NewQueries(db DBTX) *Queries {
	return &Queries{db: db}
}

// Connect to SQLite database and ensure schema exists
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/tasks.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

func WithTx(ctx context.Context, db *sql.DB, fn func(*Queries) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQueries(tx)
	if err := fn(q); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

const (
	createTaskSQL = `
	INSERT INTO tasks (id, title, start, end, status, user_id)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	updateTaskSQL = `
	UPDATE tasks
	SET title = ?, start = ?, end = ?, status = ?, user_id = ?
	WHERE id = ?
	`
	deleteTaskSQL              = `DELETE FROM tasks WHERE id = ?`
	getTaskSQL                 = `SELECT id, title, start, end, status, user_id FROM tasks WHERE id = ?`
	getTasksSQL                = `SELECT id, title, start, end, status, user_id FROM tasks WHERE user_id = ?`
	checkTaskOverlapSQL        = `SELECT 1 FROM tasks WHERE (? < end) AND (? > start)`
	checkTaskOverlapExcludeSQL = `SELECT 1 FROM tasks WHERE (? < end) AND (? > start) AND id != ?`
)

// CreateTask inserts a new task into DB
func (q *Queries) CreateTask(ctx context.Context, t models.Task) error {
	// check if no task overlaps
	result, err := q.db.QueryContext(ctx, checkTaskOverlapSQL, t.Start, t.End)
	if err != nil {
		return err
	}
	defer result.Close()

	if result.Next() {
		return ErrTaskOverlap
	}

	_, err = q.db.ExecContext(ctx, createTaskSQL, t.ID, t.Title, t.Start, t.End, t.Status, t.UserID)
	return err
}

// UpdateTask updates an existing task
func (q *Queries) UpdateTask(ctx context.Context, t models.Task) error {
	// Check for overlap with other tasks (excluding this task)
	result, err := q.db.QueryContext(ctx, checkTaskOverlapExcludeSQL, t.Start, t.End, t.ID)
	if err != nil {
		return err
	}
	defer result.Close()

	if result.Next() {
		return ErrTaskOverlap
	}

	// Update the task
	execResult, err := q.db.ExecContext(ctx, updateTaskSQL, t.Title, t.Start, t.End, t.Status, t.UserID, t.ID)
	if err != nil {
		return err
	}

	rows, err := execResult.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteTask deletes a task by ID
func (q *Queries) DeleteTask(ctx context.Context, id string) error {
	result, err := q.db.ExecContext(ctx, deleteTaskSQL, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// GetTask retrieves a task by ID
func (q *Queries) GetTask(ctx context.Context, id string) (*models.Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskSQL, id)

	var t models.Task
	err := row.Scan(&t.ID, &t.Title, &t.Start, &t.End, &t.Status, &t.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// CheckTaskOverlap checks if a task overlaps with any existing tasks
func (q *Queries) CheckTaskOverlap(ctx context.Context, start, end time.Time) error {
	result, err := q.db.QueryContext(ctx, checkTaskOverlapSQL, start, end)
	if err != nil {
		return err
	}
	defer result.Close()

	if result.Next() {
		return ErrTaskOverlap
	}
	return nil
}

// GetTasks retrieves all tasks for a user
func (q *Queries) GetTasks(ctx context.Context, userID string) ([]*models.Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Start, &t.End, &t.Status, &t.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}
