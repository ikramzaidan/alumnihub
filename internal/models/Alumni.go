package models

import "time"

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

type AlumniEducation struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	School            string    `json:"school_name"`
	Degree            string    `json:"school_degree"`
	StudyMajor        string    `json:"school_study_major"`
	StartYear         int       `json:"start_year"`
	EndYear           int       `json:"end_year"`
	CurrentlyStudying bool      `json:"currently_studying,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type AlumniJob struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	Position         string    `json:"position"`
	Company          string    `json:"company"`
	CompanyLocation  string    `json:"company_location"`
	EmploymentType   string    `json:"employment_type"`
	StartYear        int       `json:"start_year"`
	EndYear          int       `json:"end_year"`
	CurrentlyWorking bool      `json:"currently_working,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
