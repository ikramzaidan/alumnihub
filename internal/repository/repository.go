package repository

import (
	"alumnihub/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllAlumni() ([]*models.Alumni, error)
	Alumni(id int) (*models.Alumni, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)

	AllArticles() ([]*models.Article, error)
	Article(id int) (*models.Article, error)
	InsertArticle(article models.Article) (int, error)
	UpdateArticle(article models.Article) error
	DeleteArticle(id int) error

	AllForms() ([]*models.Form, error)
	Form(id int) (*models.Form, error)
	ShowForm(id int) (*models.Form, error)
	InsertForm(form models.Form) (int, error)
	UpdateForm(form models.Form) error
	DeleteForm(id int) error
	ShowFormAnswers(id int) (*models.Form, error)

	Question(id int) (*models.Question, error)
	QuestionsByForm(id int) ([]*models.Question, error)
	InsertQuestion(question models.Question) (int, error)
	UpdateQuestionOptions(id int, options []string) error

	InsertAnswers(answers []*models.Answer) error
}
