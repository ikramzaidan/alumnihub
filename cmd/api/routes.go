package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Post("/authenticate", app.authenticate)
	mux.Post("/register_check", app.registerCheck)
	mux.Post("/register", app.register)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	mux.Get("/public/{image_path}", app.serveImage)
	mux.Get("/forms/{id}/answers/export", app.exportAnswers)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/alumni", app.AllAlumni)
		mux.Get("/alumni/{id}", app.Alumni)

		mux.Get("/articles", app.allArticles)
		mux.Get("/articles/{slug}", app.article)

		mux.Get("/forms", app.allForms)                   // Get all forms data
		mux.Get("/forms/{id}", app.form)                  // Get a form data without questions
		mux.Get("/forms/{id}/show", app.showForm)         // Get a complete form data with questions and options within the form
		mux.Post("/forms/{id}/submit", app.insertAnswers) // Submit form answers

		mux.Get("/forums", app.allForums) // Get all forum data
		mux.Get("/forums/{id}", app.forum)
		mux.Get("/forums/user/{username}", app.allUserForums)
		mux.Post("/forums/create", app.insertForum)
		mux.Delete("/forums/{id}", app.deleteForum)
		mux.Post("/forums/{id}/like", app.insertLike)
		mux.Post("/forums/{id}/unlike", app.deleteLike)
		mux.Post("/forums/{id}/reply", app.insertComment)

		mux.Get("/profile", app.myProfile)
		mux.Get("/profile/{username}", app.profile)
		mux.Patch("/profile/update", app.updateProfile)
		mux.Post("/profile/educations/create", app.insertAlumniEducation)
		mux.Delete("/profile/educations/{id}", app.deleteAlumniEducation)
		mux.Post("/profile/jobs/create", app.insertAlumniJob)
		mux.Delete("/profile/jobs/{id}", app.deleteAlumniJob)

		mux.Get("/jobs", app.allJobs)
		mux.Get("/jobs/{id}", app.job)
		mux.Post("/jobs/create", app.insertJob)
		mux.Patch("/jobs/{id}", app.updateJob)
		mux.Delete("/jobs/{id}", app.deleteJob)

		mux.Get("/likes", app.userLikes)
		mux.Get("/answers", app.userAnswers)

		mux.Post("/upload_image", app.uploadImage)

		// Routes untuk admin
		mux.Route("/", func(mux chi.Router) {
			mux.Use(app.adminRequired)

			mux.Get("/dashboard", app.Dashboard)

			mux.Post("/alumni/create", app.insertAlumni)
			mux.Post("/alumni/import", app.importAlumni)
			mux.Post("/alumni/import/save", app.insertImportAlumni)
			mux.Patch("/alumni/{id}", app.updateAlumni)
			mux.Delete("/alumni/{id}", app.deleteAlumni)

			mux.Get("/articles/{id}/show", app.showArticle)
			mux.Post("/articles/create", app.insertArticle)
			mux.Patch("/articles/{id}", app.updateArticle)
			mux.Delete("/articles/{id}", app.deleteArticle)

			mux.Post("/forms/create", app.insertForm)
			mux.Patch("/forms/{id}", app.updateForm)
			mux.Delete("/forms/{id}", app.deleteForm)
			mux.Get("/forms/{id}/answers", app.showFormAnswers)
			mux.Get("/forms/{fid}/questions/{qid}/answers", app.showQuestionAnswers)

			mux.Get("/questions/{id}", app.question)
			mux.Post("/questions/create", app.insertQuestion)
			mux.Delete("/questions/{id}", app.deleteQuestion)
			mux.Patch("/questions/{id}", app.updateQuestion)
		})
	})

	return mux
}
