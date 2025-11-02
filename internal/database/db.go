package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Adjanour/vesper/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("duplicate")
var ErrInvalid = errors.New("invalid")
var ErrUnauthorized = errors.New("unauthorized")

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Queries struct {
	db DBTX
}

func NewQueries(db DBTX) *Queries {
	return &Queries{
		db: db,
	}
}


func (q *Queries) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return q.db.QueryContext(ctx, query, args...)
}

func (q *Queries) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return q.db.QueryRowContext(ctx, query, args...)
}

const createTaskSQL = `
INSERT INTO tasks (id, title, start, end, status, user_id)
VALUES (?, ?, ?, ?, ?, ?)
`

const updateTaskSQL = `
UPDATE tasks
SET title = ?, start = ?, end = ?, status = ?, user_id = ?
WHERE id = ?
`
const deleteTaskSQL = `DELETE FROM tasks WHERE id = ?`
const getTaskSQL = `SELECT id, title, start, end, status, user_id FROM tasks WHERE id = ?`

func (q *Queries) CreateTask(ctx context.Context, t models.Task) error {
	_, err := q.db.ExecContext(ctx, createTaskSQL, t.ID, t.Title, t.Start, t.End, t.Status, t.UserID)
	return err
}

func (q *Queries) UpdateTask(ctx context.Context, t models.Task) error {
	result, err := q.db.ExecContext(ctx, updateTaskSQL, t.Title, t.Start, t.End, t.Status, t.UserID, t.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}


func (q *Queries) DeleteTask(ctx context.Context, id string) error {
	result, err := q.db.ExecContext(ctx, deleteTaskSQL, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (q *Queries) GetTask(ctx context.Context, id string) (*models.Task, error) {
	row := q.QueryRow(ctx, getTaskSQL, id)
	var t models.Task
	err := row.Scan(&t.ID, &t.Title, &t.Start, &t.End, &t.Status, &t.UserID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
