package entity

import "time"

type Task struct {
	ID          int       `json:"id"`
	Done        bool      `json:"done"`
	UserID      string    `json:"user_id"`
	Header      string    `json:"header"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	StartDate   time.Time `json:"start_date"`
}
