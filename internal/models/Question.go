package models

import "time"

type Question struct {
	ID                int            `json:"id"`
	FormID            int            `json:"form_id"`
	Question          string         `json:"question_text"`
	Type              string         `json:"type"`
	Extension         bool           `json:"extension"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	Options           []*Option      `json:"options,omitempty"`
	OptionsArray      []string       `json:"options_array,omitempty"`
	Answers           []*Answer      `json:"answers,omitempty"`
	GroupAnswer       []*GroupAnswer `json:"answers_group,omitempty"`
	QuestionExtension *Extension     `json:"question_extension,omitempty"`
}

type Option struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Option     string `json:"option_text"`
}

type Extension struct {
	ID               int    `json:"id"`
	QuestionID       int    `json:"question_id"`
	FollowUpQuestion int    `json:"followup_question_id"`
	FollowUpOption   string `json:"followup_option_value"`
}
