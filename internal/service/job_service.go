package service

import (
	"errors"
	"time"

	"alumnihub/internal/models"
	"alumnihub/internal/repository"
)

type JobService struct {
	Repo repository.DatabaseRepo
}

func NewJobService(repo repository.DatabaseRepo) *JobService {
	return &JobService{Repo: repo}
}

func (s *JobService) AllJobs() ([]*models.Job, error) {
	return s.Repo.AllJobs()
}

func (s *JobService) GetJob(id int) (*models.Job, error) {
	return s.Repo.Job(id)
}

func (s *JobService) CreateJob(userID int, job models.Job) error {
	job.UserID = userID
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	_, err := s.Repo.InsertJob(job)
	return err
}

func (s *JobService) UpdateJob(jobID int, payload models.Job) error {
	if jobID != payload.ID {
		return errors.New("invalid request")
	}

	job, err := s.Repo.Job(payload.ID)
	if err != nil {
		return err
	}

	job.JobPosition = payload.JobPosition
	job.Company = payload.Company
	job.JobLocation = payload.JobLocation
	job.JobType = payload.JobType
	job.MinSalary = payload.MinSalary
	job.MaxSalary = payload.MaxSalary
	job.Description = payload.Description
	job.Closed = payload.Closed
	job.UpdatedAt = time.Now()

	return s.Repo.UpdateJob(*job)
}

func (s *JobService) DeleteJob(userID int, jobID int) error {
	if jobID <= 0 {
		return errors.New("invalid job id")
	}

	deleted, err := s.Repo.DeleteJob(userID, jobID)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.New("job not found or unauthorized")
	}

	return nil
}
