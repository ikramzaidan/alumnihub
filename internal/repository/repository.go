package repository

import (
	"alumnihub/internal/models"
	"database/sql"
)

type UserRepository interface {
	InsertUser(user models.User) (int, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserUsernameByID(id int) (string, error)
	GetUserIDByUsername(username string) (int, error)
	GetUserPhotoByID(id int) (string, error)
}

type AlumniRepository interface {
	AllAlumni() ([]*models.Alumni, error)
	Alumni(id int) (*models.Alumni, error)
	InsertAlumni(alumni models.Alumni) error
	UpdateAlumni(alumni models.Alumni) error
	DeleteAlumni(id int) error
	GetAlumniByNISN(nisn string) (*models.Alumni, error)
	GetAlumniNameByID(id int) (string, error)
}

type DashboardRepository interface {
	CountAlumni() (int, error)
	CountAlumniAccount() (int, error)
}

type ProfileRepository interface {
	GetProfiles() ([]*models.Profile, error)
	InsertProfile(profile models.Profile) (int, error)
	UpdateProfile(profile models.Profile) error
	GetProfileByAlumniID(id int) (*models.Profile, error)
	GetProfileByUserID(id int) (*models.Profile, error)
	UpdateAdminProfile(profile models.Profile) error
	GetAdminProfileByUserID(id int) (*models.Profile, error)
	InsertAlumniEducation(education models.AlumniEducation) error
	GetAlumniEducations(id int) ([]*models.AlumniEducation, error)
	GetAlumniEducation(id int) (*models.AlumniEducation, error)
	DeleteAlumniEducations(id int) error
	InsertAlumniJob(alumnijob models.AlumniJob) error
	GetAlumniJobs(id int) ([]*models.AlumniJob, error)
	GetAlumniJob(id int) (*models.AlumniJob, error)
	DeleteAlumniJobs(id int) error
}

type ArticleRepository interface {
	AllArticles() ([]*models.Article, error)
	Article(id int) (*models.Article, error)
	ArticleBySlug(slug string) (*models.Article, error)
	InsertArticle(article models.Article) (int, error)
	UpdateArticle(article models.Article) error
	DeleteArticle(id int) error
}

type FormRepository interface {
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
	UpdateQuestion(question models.Question) error
	DeleteQuestion(id int) error
	UpdateQuestionOptions(id int, options []string) error
	DeleteQuestionOptions(id int) error
	UpdateQuestionExtension(extension *models.Extension) error
	DeleteQuestionExtension(id int) error
	InsertAnswers(answers []*models.Answer) error
	GroupAnswersByQuestion(forumID int, questionID int) ([]*models.GroupAnswer, error)
	GetAnswersByUser(id int) ([]*models.Answer, error)
}

type ForumRepository interface {
	AllForums() ([]*models.Forum, error)
	Forum(id int) (*models.Forum, error)
	InsertForum(forum models.Forum) (int, error)
	DeleteForum(userId int, forumId int) (bool, error)
	GetForumsByUser(id int) ([]*models.Forum, error)
	InsertComment(comment models.Comment) (int, error)
	GetCommentsByForum(id int) ([]*models.Comment, error)
	InsertLike(like models.Like) error
	DeleteLike(userId int, forumId int) error
	GetLikesByUser(id int) ([]*models.Like, error)
	GetLikesByForum(id int) ([]*models.Like, error)
}

type JobRepository interface {
	AllJobs() ([]*models.Job, error)
	Job(id int) (*models.Job, error)
	InsertJob(job models.Job) (int, error)
	UpdateJob(job models.Job) error
	DeleteJob(id int) error
}

type DatabaseRepo interface {
	Connection() *sql.DB

	UserRepository
	AlumniRepository
	DashboardRepository
	ProfileRepository
	ArticleRepository
	FormRepository
	ForumRepository
	JobRepository
}
