package service

import (
	"alumnihub/internal/models"
	"alumnihub/internal/repository"
)

type DashboardService struct {
	Repo repository.DatabaseRepo
}

func NewDashboardService(repo repository.DatabaseRepo) *DashboardService {
	return &DashboardService{Repo: repo}
}

type DashboardData struct {
	CountAlumni        int               `json:"count_alumni"`
	CountAlumniAccount int               `json:"count_alumni_account"`
	Profiles           []*models.Profile `json:"profiles,omitempty"`
}

func (s *DashboardService) GetSummary() (*DashboardData, error) {
	countAlumni, err := s.Repo.CountAlumni()
	if err != nil {
		return nil, err
	}

	countAlumniAccount, err := s.Repo.CountAlumniAccount()
	if err != nil {
		return nil, err
	}

	profiles, err := s.Repo.GetProfiles()
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		CountAlumni:        countAlumni,
		CountAlumniAccount: countAlumniAccount,
		Profiles:           profiles,
	}, nil
}
