package models

type Alumni struct {
	ID           int    `json:"id"`
	NISN         string `json:"nisn,omitempty"`
	NIS          string `json:"nis,omitempty"`
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	Phone        string `json:"phone"`
	Year         int    `json:"graduation_year"`
	Class        string `json:"class"`
	UserUsername string `json:"user_username,omitempty"`
}
