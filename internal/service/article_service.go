package service

import (
	"errors"
	"time"

	"alumnihub/internal/models"
	"alumnihub/internal/repository"
	"alumnihub/internal/utils"
)

type ArticleService struct {
	Repo repository.DatabaseRepo
}

func NewArticleService(repo repository.DatabaseRepo) *ArticleService {
	return &ArticleService{Repo: repo}
}

func (s *ArticleService) All() ([]*models.Article, error) {
	return s.Repo.AllArticles()
}

func (s *ArticleService) GetBySlug(slug string) (*models.Article, error) {
	return s.Repo.ArticleBySlug(slug)
}

func (s *ArticleService) GetByID(articleID int) (*models.Article, error) {
	return s.Repo.Article(articleID)
}

func (s *ArticleService) Create(article models.Article) error {
	imgSrc, err := utils.GetFirstImageFromHTML(article.Body)
	if err != nil {
		return err
	}

	article.Image = imgSrc
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	article.PublishedAt = time.Now()

	_, err = s.Repo.InsertArticle(article)
	return err
}

func (s *ArticleService) Update(articleID int, payload models.Article) error {
	if articleID != payload.ID {
		return errors.New("invalid request")
	}

	article, err := s.Repo.Article(payload.ID)
	if err != nil {
		return err
	}

	imgSrc, err := utils.GetFirstImageFromHTML(payload.Body)
	if err != nil {
		return err
	}

	article.Image = imgSrc
	article.Title = payload.Title
	article.Slug = payload.Slug
	article.Body = payload.Body
	article.Status = payload.Status
	article.UpdatedAt = time.Now()
	if article.Status != "published" && payload.Status == "published" {
		article.PublishedAt = time.Now()
	}

	return s.Repo.UpdateArticle(*article)
}

func (s *ArticleService) Delete(id int) error {
	return s.Repo.DeleteArticle(id)
}
