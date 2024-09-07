package models

type Profile struct {
	ID           int                `json:"id,omitempty"`
	UserID       int                `json:"user_id,omitempty"`
	AlumniID     int                `json:"alumni_id,omitempty"`
	UserName     string             `json:"user_name,omitempty"`
	UserUsername string             `json:"user_username,omitempty"`
	Bio          string             `json:"bio,omitempty"`
	Location     string             `json:"location,omitempty"`
	Facebook     string             `json:"sm_facebook,omitempty"`
	Instagram    string             `json:"sm_instagram,omitempty"`
	Twitter      string             `json:"sm_twitter,omitempty"`
	Tiktok       string             `json:"sm_tiktok,omitempty"`
	Photo        string             `json:"photo,omitempty"`
	Educations   []*AlumniEducation `json:"educations,omitempty"`
	Jobs         []*AlumniJob       `json:"jobs,omitempty"`
}
