package models

import "time"

type Form struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Hidden         string      `json:"hidden"`
	HasTimeLimit   string      `json:"has_time_limit"`
	StartDate      time.Time   `json:"start_date"`
	EndDate        time.Time   `json:"end_date"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Questions      []*Question `json:"questions,omitempty"`
	QuestionsArray []int       `json:"questions_array,omitempty"`
}
