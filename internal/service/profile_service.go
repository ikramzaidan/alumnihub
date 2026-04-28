package service

import (
	"errors"
	"strconv"

	"alumnihub/internal/auth"
	"alumnihub/internal/models"
	"alumnihub/internal/repository"
)

type ProfileService struct {
	Repo repository.DatabaseRepo
}

func NewProfileService(repo repository.DatabaseRepo) *ProfileService {
	return &ProfileService{Repo: repo}
}

func (s *ProfileService) GetByUsername(username string) (*models.Profile, error) {
	userID, err := s.Repo.GetUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.buildProfile(userID, user.IsAdmin)
}

func (s *ProfileService) GetMyProfile(claims *auth.Claims) (*models.Profile, error) {
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	return s.buildProfile(userID, claims.IsAdmin)
}

func (s *ProfileService) UpdateProfile(claims *auth.Claims, payload models.Profile) error {
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return errors.New("invalid user ID in token")
	}

	var profile models.Profile
	if !claims.IsAdmin {
		profileAlumni, err := s.Repo.GetProfileByUserID(userID)
		if err != nil {
			return err
		}

		alumniName, err := s.Repo.GetAlumniNameByID(profileAlumni.AlumniID)
		if err != nil {
			return err
		}

		profile = *profileAlumni
		profile.UserName = alumniName
		profile.Location = payload.Location
	} else {
		profileAdmin, err := s.Repo.GetAdminProfileByUserID(userID)
		if err != nil {
			return err
		}
		profile = *profileAdmin
	}

	profile.Bio = payload.Bio
	profile.Facebook = payload.Facebook
	profile.Instagram = payload.Instagram
	profile.Twitter = payload.Twitter
	profile.Tiktok = payload.Tiktok
	profile.Photo = payload.Photo

	if !claims.IsAdmin {
		return s.Repo.UpdateProfile(profile)
	}
	return s.Repo.UpdateAdminProfile(profile)
}

func (s *ProfileService) AddEducation(userID int, education models.AlumniEducation) error {
	education.UserID = userID
	return s.Repo.InsertAlumniEducation(education)
}

func (s *ProfileService) DeleteEducation(userID, educationID int) error {
	education, err := s.Repo.GetAlumniEducation(educationID)
	if err != nil {
		return err
	}

	if userID != education.UserID {
		return errors.New("user have no permissions")
	}

	return s.Repo.DeleteAlumniEducations(educationID)
}

func (s *ProfileService) AddJob(userID int, alumnijob models.AlumniJob) error {
	alumnijob.UserID = userID
	return s.Repo.InsertAlumniJob(alumnijob)
}

func (s *ProfileService) DeleteJob(userID, jobID int) error {
	alumnijob, err := s.Repo.GetAlumniJob(jobID)
	if err != nil {
		return err
	}

	if userID != alumnijob.UserID {
		return errors.New("user have no permissions")
	}

	return s.Repo.DeleteAlumniJobs(jobID)
}

func (s *ProfileService) buildProfile(userID int, isAdmin bool) (*models.Profile, error) {
	var profile models.Profile
	if !isAdmin {
		profileAlumni, err := s.Repo.GetProfileByUserID(userID)
		if err != nil {
			return nil, err
		}

		alumniName, err := s.Repo.GetAlumniNameByID(profileAlumni.AlumniID)
		if err != nil {
			return nil, err
		}

		profile = *profileAlumni
		profile.UserName = alumniName
	} else {
		profileAdmin, err := s.Repo.GetAdminProfileByUserID(userID)
		if err != nil {
			return nil, err
		}
		profile = *profileAdmin
	}

	userUsername, err := s.Repo.GetUserUsernameByID(userID)
	if err != nil {
		return nil, err
	}

	userPhoto, err := s.Repo.GetUserPhotoByID(userID)
	if err != nil {
		return nil, err
	}

	educations, err := s.Repo.GetAlumniEducations(userID)
	if err != nil {
		return nil, err
	}

	jobs, err := s.Repo.GetAlumniJobs(userID)
	if err != nil {
		return nil, err
	}

	profile.UserUsername = userUsername
	profile.Photo = userPhoto
	profile.Educations = educations
	profile.Jobs = jobs

	return &profile, nil
}
