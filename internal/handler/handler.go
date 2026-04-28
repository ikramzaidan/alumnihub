package handler

import (
	"alumnihub/internal/auth"
	"alumnihub/internal/service"
)

type Handler struct {
	AuthService      *service.AuthService
	ProfileService   *service.ProfileService
	AlumniService    *service.AlumniService
	ArticleService   *service.ArticleService
	FormService      *service.FormService
	ForumService     *service.ForumService
	JobService       *service.JobService
	DashboardService *service.DashboardService
	Auth             *auth.Auth
}

func NewHandler(
	authService *service.AuthService,
	profileService *service.ProfileService,
	alumniService *service.AlumniService,
	articleService *service.ArticleService,
	formService *service.FormService,
	forumService *service.ForumService,
	jobService *service.JobService,
	dashboardService *service.DashboardService,
	auth *auth.Auth,
) *Handler {
	return &Handler{
		AuthService:      authService,
		ProfileService:   profileService,
		AlumniService:    alumniService,
		ArticleService:   articleService,
		FormService:      formService,
		ForumService:     forumService,
		JobService:       jobService,
		DashboardService: dashboardService,
		Auth:             auth,
	}
}
