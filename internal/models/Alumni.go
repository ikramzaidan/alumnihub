package models

type Alumni struct {
	ID     int    `json:"id"`
	NISN   string `json:"nisn,omitempty"`
	NIS    string `json:"nis,omitempty"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
	Year   int    `json:"graduation_year"`
	Class  string `json:"class"`
}

type Profile struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id,omitempty"`
	AlumniID  int    `json:"alumni_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Location  string `json:"location,omitempty"`
	Facebook  string `json:"sm_facebook,omitempty"`
	Instagram string `json:"sm_instagram,omitempty"`
	Twitter   string `json:"sm_twitter,omitempty"`
	Tiktok    string `json:"sm_tiktok,omitempty"`
}
