package main

import (
	"alumnihub/internal/auth"
	"alumnihub/internal/handler"
	internalMiddleware "alumnihub/internal/middleware"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func routes(handlers *handler.Handler, authSvc *auth.Auth) http.Handler {
	mux := chi.NewRouter()

	mux.Use(chiMiddleware.Recoverer)
	mux.Use(internalMiddleware.EnableCORS)

	mux.Get("/", handlers.Home)
	mux.Post("/authenticate", handlers.Authenticate)
	mux.Post("/register_check", handlers.RegisterCheck)
	mux.Post("/register", handlers.Register)
	mux.Post("/forgot_password", handlers.ForgotPassword)
	mux.Post("/reset_password", handlers.ResetPassword)
	mux.Get("/refresh", handlers.RefreshToken)
	mux.Get("/logout", handlers.Logout)
	mux.Get("/public/{image_path}", handlers.ServeImage)
	mux.Get("/forms/{id}/answers/export", handlers.ExportAnswers)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(internalMiddleware.AuthRequired(authSvc))

		mux.Get("/alumni", handlers.AllAlumni)
		mux.Get("/alumni/{id}", handlers.Alumni)

		mux.Get("/articles", handlers.AllArticles)
		mux.Get("/articles/{slug}", handlers.Article)

		mux.Get("/forms", handlers.AllForms)
		mux.Get("/forms/{id}", handlers.Form)
		mux.Get("/forms/{id}/show", handlers.ShowForm)
		mux.Post("/forms/{id}/submit", handlers.InsertAnswers)

		mux.Get("/forums", handlers.AllForums)
		mux.Get("/forums/{id}", handlers.Forum)
		mux.Get("/forums/user/{username}", handlers.AllUserForums)
		mux.Post("/forums", handlers.InsertForum)
		mux.Delete("/forums/{id}", handlers.DeleteForum)
		mux.Post("/forums/{id}/like", handlers.InsertLike)
		mux.Post("/forums/{id}/unlike", handlers.DeleteLike)
		mux.Post("/forums/{id}/reply", handlers.InsertComment)

		mux.Get("/profile", handlers.MyProfile)
		mux.Get("/profile/{username}", handlers.Profile)
		mux.Patch("/profile", handlers.UpdateProfile)
		mux.Post("/profile/educations", handlers.InsertAlumniEducation)
		mux.Delete("/profile/educations/{id}", handlers.DeleteAlumniEducation)
		mux.Post("/profile/jobs", handlers.InsertAlumniJob)
		mux.Delete("/profile/jobs/{id}", handlers.DeleteAlumniJob)

		mux.Get("/jobs", handlers.AllJobs)
		mux.Get("/jobs/{id}", handlers.Job)
		mux.Post("/jobs", handlers.InsertJob)
		mux.Patch("/jobs/{id}", handlers.UpdateJob)
		mux.Delete("/jobs/{id}", handlers.DeleteJob)

		mux.Get("/likes", handlers.UserLikes)
		mux.Get("/answers", handlers.UserAnswers)

		mux.Post("/upload_image", handlers.UploadImage)

		mux.Route("/", func(mux chi.Router) {
			mux.Use(internalMiddleware.AdminRequired(authSvc))

			mux.Get("/dashboard", handlers.Dashboard)

			mux.Post("/alumni", handlers.InsertAlumni)
			mux.Post("/alumni/import", handlers.ImportAlumni)
			mux.Post("/alumni/import/save", handlers.InsertImportAlumni)
			mux.Patch("/alumni/{id}", handlers.UpdateAlumni)
			mux.Delete("/alumni/{id}", handlers.DeleteAlumni)

			mux.Get("/articles/{id}/show", handlers.ShowArticle)
			mux.Post("/articles", handlers.InsertArticle)
			mux.Patch("/articles/{id}", handlers.UpdateArticle)
			mux.Delete("/articles/{id}", handlers.DeleteArticle)

			mux.Post("/forms", handlers.InsertForm)
			mux.Patch("/forms/{id}", handlers.UpdateForm)
			mux.Delete("/forms/{id}", handlers.DeleteForm)
			mux.Get("/forms/{id}/answers", handlers.ShowFormAnswers)
			mux.Get("/forms/{fid}/questions/{qid}/answers", handlers.ShowQuestionAnswers)

			mux.Get("/questions/{id}", handlers.Question)
			mux.Post("/questions", handlers.InsertQuestion)
			mux.Delete("/questions/{id}", handlers.DeleteQuestion)
			mux.Patch("/questions/{id}", handlers.UpdateQuestion)
		})
	})

	return mux
}
