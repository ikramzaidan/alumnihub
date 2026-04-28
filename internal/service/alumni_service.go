package service

import (
	"alumnihub/internal/models"
	"alumnihub/internal/repository"
	"database/sql"
)

type AlumniService struct {
	Repo repository.DatabaseRepo
}

func NewAlumniService(repo repository.DatabaseRepo) *AlumniService {
	return &AlumniService{Repo: repo}
}

func (s *AlumniService) All() ([]*models.Alumni, error) {
	return s.Repo.AllAlumni()
}

func (s *AlumniService) Get(alumniID int) (*models.Alumni, error) {
	alumni, err := s.Repo.Alumni(alumniID)
	if err != nil {
		return nil, err
	}

	profile, err := s.Repo.GetProfileByAlumniID(alumniID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if profile != nil {
		user, err := s.Repo.GetUserByID(profile.UserID)
		if err != nil {
			return nil, err
		}
		alumni.UserUsername = user.Username
	}

	return alumni, nil
}

func (s *AlumniService) Create(alumni models.Alumni) error {
	return s.Repo.InsertAlumni(alumni)
}

func (s *AlumniService) InsertBulk(alumniList []models.Alumni) error {
	for _, alumni := range alumniList {
		if err := s.Repo.InsertAlumni(alumni); err != nil {
			return err
		}
	}
	return nil
}

func (s *AlumniService) Update(alumniID int, payload models.Alumni) error {
	alumni, err := s.Repo.Alumni(alumniID)
	if err != nil {
		return err
	}

	alumni.Name = payload.Name
	alumni.Gender = payload.Gender
	alumni.Phone = payload.Phone
	alumni.Year = payload.Year
	alumni.Class = payload.Class
	alumni.NISN = payload.NISN
	alumni.NIS = payload.NIS

	return s.Repo.UpdateAlumni(*alumni)
}

func (s *AlumniService) Delete(id int) error {
	return s.Repo.DeleteAlumni(id)
}
