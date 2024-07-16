package models

import (
	"time"
)

type Job struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	JobPosition string    `json:"job_position,omitempty"`
	Company     string    `json:"company,omitempty"`
	JobLocation string    `json:"job_location,omitempty"`
	JobType     string    `json:"job_type,omitempty"`
	MinSalary   int       `json:"min_salary,omitempty"`
	MaxSalary   int       `json:"max_salary,omitempty"`
	Description string    `json:"description,omitempty"`
	Closed      bool      `json:"closed,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
