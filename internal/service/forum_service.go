package service

import (
	"errors"
	"time"

	"alumnihub/internal/models"
	"alumnihub/internal/repository"
)

type ForumService struct {
	Repo repository.DatabaseRepo
}

func NewForumService(repo repository.DatabaseRepo) *ForumService {
	return &ForumService{Repo: repo}
}

func (s *ForumService) AllForums() ([]*models.Forum, error) {
	return s.Repo.AllForums()
}

func (s *ForumService) GetForumsByUser(username string) ([]*models.Forum, error) {
	userID, err := s.Repo.GetUserIDByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.Repo.GetForumsByUser(userID)
}

func (s *ForumService) GetForum(id int) (*models.Forum, error) {
	return s.Repo.Forum(id)
}

func (s *ForumService) CreateForum(userID int, forumText string) error {
	forum := models.Forum{
		Forum:       forumText,
		UserID:      userID,
		PublishedAt: time.Now(),
	}
	_, err := s.Repo.InsertForum(forum)
	return err
}

func (s *ForumService) DeleteForum(userID, forumID int) error {
	if forumID <= 0 {
		return errors.New("invalid forum id")
	}

	deleted, err := s.Repo.DeleteForum(userID, forumID)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.New("forum not found or unauhthorized")
	}
	return nil
}

func (s *ForumService) InsertComment(userID, forumID int, commentText string) error {
	comment := models.Comment{
		Comment:     commentText,
		UserID:      userID,
		ForumID:     forumID,
		PublishedAt: time.Now(),
	}
	_, err := s.Repo.InsertComment(comment)
	return err
}

func (s *ForumService) GetLikesByUser(userID int) ([]*models.Like, error) {
	return s.Repo.GetLikesByUser(userID)
}

func (s *ForumService) AddLike(userID, forumID int) error {
	like := models.Like{
		ForumID:   forumID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return s.Repo.InsertLike(like)
}

func (s *ForumService) RemoveLike(userID, forumID int) error {
	return s.Repo.DeleteLike(userID, forumID)
}
