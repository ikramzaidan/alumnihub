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

	mux.Route("/", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/alumni", app.AllAlumni)
		mux.Get("/alumni/{id}", app.Alumni)

		mux.Get("/articles", app.allArticles)
		mux.Get("/articles/{id}", app.article)

		mux.Get("/forms", app.allForms)                   // Get all forms data
		mux.Get("/forms/{id}", app.form)                  // Get a form data without questions
		mux.Get("/forms/{id}/show", app.showForm)         // Get a complete form data with questions and options within the form
		mux.Post("/forms/{id}/submit", app.insertAnswers) // Submit form answers

		mux.Post("/upload_image", app.uploadImage)

		// Routes untuk admin
		mux.Route("/", func(mux chi.Router) {
			mux.Use(app.adminRequired)

			mux.Post("/alumni/create", app.insertAlumni)
			mux.Patch("/alumni/{id}", app.updateAlumni)
			mux.Delete("/alumni/{id}", app.deleteAlumni)

			mux.Post("/articles/create", app.insertArticle)
			mux.Patch("/articles/{id}", app.updateArticle)
			mux.Delete("/articles/{id}", app.deleteArticle)

			mux.Post("/forms/create", app.insertForm)
			mux.Patch("/forms/{id}", app.updateForm)
			mux.Delete("/forms/{id}", app.deleteForm)
			mux.Get("/forms/{id}/answers", app.showFormAnswers)

			mux.Get("/questions/{id}", app.question)
			mux.Post("/questions/create", app.insertQuestion)
			mux.Patch("/questions/{id}", app.updateQuestion)
		})
	})

	return mux
}
