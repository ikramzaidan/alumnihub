package models

import "time"

type Alumni struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	BirthDate  time.Time `json:"date_of_birth"`
	BirthPlace string    `json:"place_of_birth"`
	Gender     string    `json:"gender"`
	Phone      string    `json:"phone"`
}
