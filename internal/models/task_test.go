package models

import (
	"testing"
	"time"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		status TaskStatus
		want   bool
	}{
		{StatusScheduled, true},
		{StatusDeleted, true},
		{StatusReplaced, true},
		{"invalid", false},
		{"", false},
		{"SCHEDULED", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := IsValidStatus(tt.status); got != tt.want {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestIsOverlapping(t *testing.T) {
	baseTime := time.Date(2026, 2, 8, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		newTask  Task
		existing []Task
		want     bool
	}{
		{
			name: "no overlap - before",
			newTask: Task{
				ID:     "new",
				Start:  baseTime,
				End:    baseTime.Add(1 * time.Hour),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime.Add(2 * time.Hour),
					End:    baseTime.Add(3 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: false,
		},
		{
			name: "no overlap - after",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(2 * time.Hour),
				End:    baseTime.Add(3 * time.Hour),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: false,
		},
		{
			name: "overlap - starts during",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(30 * time.Minute),
				End:    baseTime.Add(90 * time.Minute),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: true,
		},
		{
			name: "overlap - ends during",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(-30 * time.Minute),
				End:    baseTime.Add(30 * time.Minute),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: true,
		},
		{
			name: "overlap - completely contains",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(-30 * time.Minute),
				End:    baseTime.Add(90 * time.Minute),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: true,
		},
		{
			name: "overlap - completely contained",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(15 * time.Minute),
				End:    baseTime.Add(45 * time.Minute),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: true,
		},
		{
			name: "no overlap - deleted task ignored",
			newTask: Task{
				ID:     "new",
				Start:  baseTime,
				End:    baseTime.Add(1 * time.Hour),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusDeleted,
				},
			},
			want: false,
		},
		{
			name: "no overlap - replaced task ignored",
			newTask: Task{
				ID:     "new",
				Start:  baseTime,
				End:    baseTime.Add(1 * time.Hour),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusReplaced,
				},
			},
			want: false,
		},
		{
			name: "overlap - with multiple tasks",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(90 * time.Minute),
				End:    baseTime.Add(150 * time.Minute),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "task1",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
				{
					ID:     "task2",
					Start:  baseTime.Add(2 * time.Hour),
					End:    baseTime.Add(3 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: true,
		},
		{
			name: "no overlap - adjacent tasks",
			newTask: Task{
				ID:     "new",
				Start:  baseTime.Add(1 * time.Hour),
				End:    baseTime.Add(2 * time.Hour),
				Status: StatusScheduled,
			},
			existing: []Task{
				{
					ID:     "existing",
					Start:  baseTime,
					End:    baseTime.Add(1 * time.Hour),
					Status: StatusScheduled,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOverlapping(tt.newTask, tt.existing); got != tt.want {
				t.Errorf("IsOverlapping() = %v, want %v", got, tt.want)
			}
		})
	}
}
