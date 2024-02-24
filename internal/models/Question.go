package models

import "time"

type Question struct {
	ID           int       `json:"id"`
	FormID       int       `json:"form_id"`
	Question     string    `json:"question_text"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Options      []*Option `json:"options,omitempty"`
	OptionsArray []string  `json:"options_array,omitempty"`
}

type Option struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Option     string `json:"option_text"`
}
