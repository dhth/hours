package domain

import "time"

type Task struct {
	ID        int
	Summary   string
	CreatedAt time.Time
	UpdatedAt time.Time
	SecsSpent int
	Active    bool
}
