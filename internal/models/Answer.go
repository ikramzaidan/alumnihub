package models

type Answer struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	FormID     int    `json:"form_id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer_text"`
}

type GroupAnswer struct {
	FormID     int    `json:"form_id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer_text"`
	Count      int    `json:"answer_count"`
}
