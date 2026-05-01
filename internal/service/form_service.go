package service

import (
	"errors"
	"time"

	"alumnihub/internal/models"
	"alumnihub/internal/repository"
)

type FormService struct {
	Repo repository.DatabaseRepo
}

func NewFormService(repo repository.DatabaseRepo) *FormService {
	return &FormService{Repo: repo}
}

func (s *FormService) AllForms() ([]*models.Form, error) {
	return s.Repo.AllForms()
}

func (s *FormService) GetForm(formID int) (*models.Form, error) {
	return s.Repo.Form(formID)
}

func (s *FormService) ShowForm(formID int) (*models.Form, error) {
	return s.Repo.ShowForm(formID)
}

func (s *FormService) CreateForm(form models.Form) (int, error) {
	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()
	return s.Repo.InsertForm(form)
}

func (s *FormService) UpdateForm(formID int, payload models.Form) error {
	if formID != payload.ID {
		return errors.New("invalid request")
	}

	form, err := s.Repo.Form(payload.ID)
	if err != nil {
		return err
	}

	form.Title = payload.Title
	form.Description = payload.Description
	form.Hidden = payload.Hidden
	form.HasTimeLimit = payload.HasTimeLimit
	form.StartDate = payload.StartDate
	form.UpdatedAt = time.Now()

	return s.Repo.UpdateForm(*form)
}

func (s *FormService) DeleteForm(id int) error {
	return s.Repo.DeleteForm(id)
}

func (s *FormService) GetQuestion(questionID int) (*models.Question, error) {
	question, err := s.Repo.Question(questionID)
	if err != nil {
		return nil, err
	}
	groupAnswers, err := s.Repo.GroupAnswersByQuestion(question.FormID, questionID)
	if err != nil {
		return nil, err
	}
	question.GroupAnswer = groupAnswers
	return question, nil
}

func (s *FormService) CreateQuestion(question models.Question) error {
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	newID, err := s.Repo.InsertQuestion(question)
	if err != nil {
		return err
	}

	if question.Type == "multiple_choice" {
		return s.Repo.UpdateQuestionOptions(newID, question.OptionsArray)
	}

	return nil
}

func (s *FormService) UpdateQuestion(questionID int, payload models.Question) error {
	if questionID != payload.ID {
		return errors.New("invalid request")
	}

	question, err := s.Repo.Question(payload.ID)
	if err != nil {
		return err
	}

	question.Question = payload.Question
	question.Type = payload.Type
	question.Extension = payload.Extension
	question.UpdatedAt = time.Now()
	question.ID = payload.ID
	question.OptionsArray = payload.OptionsArray

	if err := s.Repo.UpdateQuestion(*question); err != nil {
		return err
	}

	if question.Type == "multiple_choice" {
		if err := s.Repo.UpdateQuestionOptions(payload.ID, question.OptionsArray); err != nil {
			return err
		}
	} else {
		if err := s.Repo.DeleteQuestionOptions(payload.ID); err != nil {
			return err
		}
	}

	if question.Extension && payload.QuestionExtension != nil {
		return s.Repo.UpdateQuestionExtension(payload.QuestionExtension)
	}

	return s.Repo.DeleteQuestionExtension(payload.ID)
}

func (s *FormService) DeleteQuestion(id int) error {
	return s.Repo.DeleteQuestion(id)
}

func (s *FormService) ShowFormAnswers(formID int) (*models.Form, error) {
	return s.Repo.ShowFormAnswers(formID)
}

func (s *FormService) InsertAnswers(answers []*models.Answer) error {
	return s.Repo.InsertAnswers(answers)
}

func (s *FormService) GetAnswersByUser(userID int) ([]*models.Answer, error) {
	return s.Repo.GetAnswersByUser(userID)
}

func (s *FormService) ShowQuestionAnswers(formID, questionID int) ([]*models.GroupAnswer, error) {
	return s.Repo.GroupAnswersByQuestion(formID, questionID)
}
