package models

import (
	"time"
)

type TaskStatus string

const (
	StatusScheduled TaskStatus = "scheduled"
	StatusDeleted   TaskStatus = "deleted"
	StatusReplaced  TaskStatus = "replaced"
)

type Task struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Start      time.Time  `json:"start"`
	End        time.Time  `json:"end"`
	UserID     string     `json:"user_id"`
	Status TaskStatus `json:"status"`
}

func IsValidStatus(s TaskStatus) bool {
	switch s {
	case StatusScheduled, StatusDeleted, StatusReplaced:
		return true
	default:
		return false
	}
}

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
}

func IsOverlapping(newTask Task, existing []Task) bool {
	for _, t := range existing {
		if t.Status == StatusDeleted || t.Status == StatusReplaced {
			continue // skip inactive blocks
		}
		if newTask.Start.Before(t.End) && newTask.End.After(t.Start) {
			return true
		}
	}
	return false
}
