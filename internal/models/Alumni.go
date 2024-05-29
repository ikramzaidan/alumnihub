package models

import "database/sql"

type Alumni struct {
	ID     int           `json:"id"`
	NISN   string        `json:"nisn,omitempty"`
	NIS    string        `json:"nis,omitempty"`
	Name   string        `json:"name"`
	Gender string        `json:"gender"`
	Phone  string        `json:"phone"`
	Year   int           `json:"graduation_year"`
	Class  string        `json:"class"`
	UserID sql.NullInt64 `json:"user_id,omitempty"`
}

func (a *Alumni) isUserIDNull(userID sql.NullInt64) bool {
	return !userID.Valid
}
